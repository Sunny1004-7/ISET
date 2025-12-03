package main

import (
	"backend/lib"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Config struct {
	lib.DatabaseConfig
	Port string
}

type Graph struct {
	Id   uint32 `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
}

type Node struct {
	Id      uint32 `gorm:"primary_key" json:"id"`
	GraphId uint32 `json:"graph_id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Src     string `json:"src"`
}

type Link struct {
	Id      uint32 `gorm:"primary_key" json:"id"`
	GraphId uint32 `json:"graph_id"`
	Source  uint32 `json:"source"`
	Target  uint32 `json:"target"`
	Type    string `json:"type"`
}

type DuplicateNode struct {
	Id    uint32 `gorm:"primary_key" json:"id"`
	Node1 uint32 `json:"node_1"`
	Node2 uint32 `json:"node_2"`
}

type ExeData struct {
	ExerciseType     string `json:"exerciseType"`
	ExerciseQuantity string `json:"exerciseQuantity"`
	SelectedModel    string `json:"selectedModel"`
	SelectedAgents   string `json:"selectedAgents"`
	Difficulty       string `json:"difficulty"`
	KnowledgePoint   string `json:"knowledgePoint"`
	Language         string `json:"language"`
	Subject          string `json:"subject"`
}

var (
	historyChannel = make(chan string, 100)
	clients        = make(map[*websocket.Conn]bool)
	clientsMutex   = sync.Mutex{}
)
var (
	// 增加写超时防止阻塞
	upgrader = websocket.Upgrader{
		HandshakeTimeout: 10 * time.Second,
		WriteBufferPool:  &sync.Pool{},
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
)

// 学生端处理习题算法选择的函数
func handleAlgorithm(exeData ExeData) string {
	log.Printf("学生端开始处理算法请求 | 题型:%s 数量:%s 模型:%s 智能体:%s",
		exeData.ExerciseType, exeData.ExerciseQuantity,
		exeData.SelectedModel, exeData.SelectedAgents)

	selectedAgentsArray := strings.Split(exeData.SelectedAgents, ",")
	for i := range selectedAgentsArray {
		selectedAgentsArray[i] = strings.TrimSpace(selectedAgentsArray[i])
	}

	// 算法路由逻辑
	switch {
	case compareSlices(selectedAgentsArray, []string{"KnowledgeAwareAgent", "DiscriminatorAgent", "SuperviseAgent", "ExeAgent"}):
		log.Printf("调用算法组合1（全智能体）")
		return runAllAgentAlgorithm(exeData)

	case compareSlices(selectedAgentsArray, []string{"KnowledgeAwareAgent", "SuperviseAgent", "ExeAgent"}):
		log.Printf("调用算法组合2（无质量评估）")
		return runNo_DisAgentAlgorithm(exeData)

	case compareSlices(selectedAgentsArray, []string{"DiscriminatorAgent", "SuperviseAgent", "ExeAgent"}):
		log.Printf("调用算法组合3（无知识感知）")
		return runNo_ktAgentAlgorithm(exeData)
	case compareSlices(selectedAgentsArray, []string{"KnowledgeAwareAgent", "DiscriminatorAgent", "SuperviseAgent"}):
		log.Printf("调用算法组合4（无习题生成）")
		return runNo_genAgentAlgorithm(exeData)
	default:
		log.Printf(" 无匹配算法 | 接收智能体: %v", selectedAgentsArray)
		return ""
	}
}

// 教师端处理习题算法选择的函数
func handleAlgorithm_tea(exeData ExeData) string {
	log.Printf("教师端开始处理算法请求 | 题型:%s 数量:%s 模型:%s 智能体:%s",
		exeData.ExerciseType, exeData.ExerciseQuantity,
		exeData.SelectedModel, exeData.SelectedAgents)

	selectedAgentsArray := strings.Split(exeData.SelectedAgents, ",")
	for i := range selectedAgentsArray {
		selectedAgentsArray[i] = strings.TrimSpace(selectedAgentsArray[i])
	}

	// 算法路由逻辑
	switch {
	case compareSlices(selectedAgentsArray, []string{"DiscriminatorAgent", "SuperviseAgent", "ExeAgent"}):
		log.Printf("调用算法组合1（全智能体）")
		return runAllAgentAlgorithm_tea(exeData)

	case compareSlices(selectedAgentsArray, []string{"SuperviseAgent", "ExeAgent"}):
		log.Printf("调用算法组合2（无质量评估）")
		return runNo_DisAgentAlgorithm_tea(exeData)
	case compareSlices(selectedAgentsArray, []string{"DiscriminatorAgent", "SuperviseAgent"}):
		log.Printf("调用算法组合3（无习题生成）")
		return runNo_genAgentAlgorithm_tea(exeData)
	default:
		log.Printf(" 无匹配算法 | 接收智能体: %v", selectedAgentsArray)
		return ""
	}
}

// 比较两个切片（不考虑顺序）
func compareSlices(slice1, slice2 []string) bool {
	
	if len(slice1) != len(slice2) {
		return false
	}
	sort.Strings(slice1)
	sort.Strings(slice2)
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}

// 学生端第一种算法全智能体
func runAllAgentAlgorithm(exeData ExeData) string {

	err := runAlgorithmScript(exeData, "../exeGen/process.py")
	if err != nil {
		log.Printf("执行失败: %v", err)
		return ""
	}

	filePath := "../exeGen/txtfile/result.txt"
	jsonStr, err := extractProcessedData(filePath)
	if err != nil {
		log.Printf("读取JSON文件失败: %v", err)
		return ""
	}
	return jsonStr
}

// 学生端第二种算法 无质量评估
func runNo_DisAgentAlgorithm(exeData ExeData) string {
	err := runAlgorithmScript(exeData, "../exeGen/process_no_discriminator.py")
	if err != nil {
		log.Printf("执行失败: %v", err)
		return ""
	}

	filePath := "../exeGen/txtfile/result_no_discriminator.txt"
	jsonStr, err := extractProcessedData(filePath)
	if err != nil {
		log.Printf("读取JSON文件失败: %v", err)
		return ""
	}
	return jsonStr
}

// 学生端第三种算法 无知识感知
func runNo_ktAgentAlgorithm(exeData ExeData) string {
	log.Printf("执行Python算法模块 | 类型:%s 数量:%s",
		exeData.ExerciseType, exeData.ExerciseQuantity)

	err := runAlgorithmScript(exeData, "../exeGen/process_no_kt.py")
	if err != nil {
		log.Printf("执行失败: %v", err)
		return ""
	}

	filePath := "../exeGen/txtfile/result_no_kt.txt"
	jsonStr, err := extractProcessedData(filePath)
	if err != nil {
		log.Printf("读取JSON文件失败: %v", err)
		return ""
	}
	
	return jsonStr
}

// 学生端第四种算法 无习题生成
func runNo_genAgentAlgorithm(exeData ExeData) string {
	log.Printf("执行Python算法模块 | 类型:%s 数量:%s",
		exeData.ExerciseType, exeData.ExerciseQuantity)

	err := runAlgorithmScript(exeData, "../exeGen/process_no_generator.py")
	if err != nil {
		log.Printf("执行失败: %v", err)
		return ""
	}
	filePath := "../exeGen/txtfile/result_no_generator.txt"
	jsonStr, err := extractProcessedData(filePath)
	if err != nil {
		log.Printf("读取JSON文件失败: %v", err)
		return ""
	}

	return jsonStr
}

// 教师端第一种算法全智能体
func runAllAgentAlgorithm_tea(exeData ExeData) string {
	//调用通用的算法执行函数
	err := runAlgorithmScript_tea(exeData, "../exeGen/tea_process.py")
	if err != nil {
		log.Printf("执行失败: %v", err)
		return ""
	}

	filePath := "../exeGen/txtfile/tea_result.txt"
	jsonStr, err := extractProcessedData(filePath)
	if err != nil {
		log.Printf("读取JSON文件失败: %v", err)
		return ""
	}

	return jsonStr
}

// 教师端端第二种算法 无质量评估
func runNo_DisAgentAlgorithm_tea(exeData ExeData) string {

	err := runAlgorithmScript_tea(exeData, "../exeGen/tea_process_no_discriminator.py")
	if err != nil {
		log.Printf("执行失败: %v", err)
		return ""
	}

	filePath := "../exeGen/txtfile/tea_result_no_discriminator.txt"
	jsonStr, err := extractProcessedData(filePath)
	if err != nil {
		log.Printf("读取JSON文件失败: %v", err)
		return ""
	}

	return jsonStr
}

// 教师端第三种算法 无习题生成
func runNo_genAgentAlgorithm_tea(exeData ExeData) string {
	log.Printf("执行Python算法模块 | 类型:%s 数量:%s",
		exeData.ExerciseType, exeData.ExerciseQuantity)

	err := runAlgorithmScript_tea(exeData, "../exeGen/tea_process_no_generator.py")
	if err != nil {
		log.Printf("执行失败: %v", err)
		return ""
	}

	filePath := "../exeGen/txtfile/tea_tea_result_no_generator.txt"
	jsonStr, err := extractProcessedData(filePath)
	if err != nil {
		log.Printf("读取JSON文件失败: %v", err)
		return ""
	}
	return jsonStr
}

// 核心处理函数
func extractProcessedData(filePath string) (string, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("文件读取失败: %v", err)
	}
	content := string(data)

	exercisesStr, err := extractLastArray(content, "exercises")
	if err != nil {
		return "", fmt.Errorf("exercises提取失败: %v", err)
	}

	chatHistoryStr, err := extractLastArray(content, "chat_history")
	if err != nil {
		return "", fmt.Errorf("chat_history提取失败: %v", err)
	}
	filteredChat, err := filterChatHistory(chatHistoryStr)
	if err != nil {
		return "", err
	}
	result := map[string]interface{}{
		"exercises":   json.RawMessage(exercisesStr),
		"chatHistory": filteredChat,
	}
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON生成失败: %v", err)
	}

	return string(jsonData), nil
}

// 通用数组提取函数
func extractLastArray(content, field string) (string, error) {
	fieldMarker := `"` + field + `"`
	startIdx := strings.LastIndex(content, fieldMarker)
	if startIdx == -1 {
		return "", fmt.Errorf("字段 %s 未找到", field)
	}
	prefix := content[startIdx:]
	arrayStart := strings.Index(prefix, "[")
	if arrayStart == -1 {
		return "", fmt.Errorf("%s 数组开始符未找到", field)
	}
	start := startIdx + arrayStart
	bracketCount := 0
	end := -1
	for i := start; i < len(content); i++ {
		switch content[i] {
		case '[':
			bracketCount++
		case ']':
			bracketCount--
			if bracketCount == 0 {
				end = i
				break
			}
		}
		if end != -1 {
			break
		}
	}
	if end == -1 {
		return "", fmt.Errorf("%s 数组结束符未找到", field)
	}

	return content[start : end+1], nil
}

// chat_history过滤逻辑
func filterChatHistory(chatStr string) ([]interface{}, error) {
	// 解析原始数据
	var rawChat []map[string]interface{}
	if err := json.Unmarshal([]byte(chatStr), &rawChat); err != nil {
		return nil, fmt.Errorf("chat_history解析失败: %v", err)
	}

	// 定义允许的agent类型
	allowed := map[string]bool{
		"agent_kt":                         true,
		"agent_host":                       true,
		"agent_generator":                  true,
		"Linguistic_Fluency_discriminator": true,
		"Knowledge_Concept_Coverage_discriminator":     true,
		"Correctness_and_Reasonableness_discriminator": true,
	}

	// 执行过滤
	filtered := make([]interface{}, 0)
	for _, item := range rawChat {
		if name, ok := item["name"].(string); ok {
			if name != "chat_manager" && allowed[name] {
				filtered = append(filtered, item)
			}
		}
	}
	return filtered, nil
}

// 学生端运行算法并将交互历史推送到 historyChannel
func runAlgorithmScript(exeData ExeData, scriptPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	// 从完整路径提取脚本名
	scriptName := filepath.Base(scriptPath)

	cmd := exec.CommandContext(
		ctx,
		"venv\\Scripts\\python.exe", 
		"-u",
		scriptName,
		"--type_of_prompt", "json_text",
		"--exercise_type", exeData.ExerciseType,
		"--output_type", "json",
		"--Number_of_Generations", exeData.ExerciseQuantity,
	)
	cmd.Dir = "../exeGen"

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("创建stdout管道失败: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("创建stderr管道失败: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动失败: %v", err)
	}

	// 实时扫描输出
	scanAndStream := func(pipe io.Reader, prefix string) {
		reader := bufio.NewReader(pipe)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					if err != io.EOF {
						log.Printf("%s读取错误: %v", prefix, err)
					}
					return
				}
				msg := strings.TrimSuffix(line, "\n")
				if msg != "" {
					
					trimmedMsg := strings.Replace(msg, prefix, "", 1)
					historyChannel <- trimmedMsg 
				}
			}
		}
	}

	// 启动输出处理协程
	go scanAndStream(stdoutPipe, "STDOUT")
	go scanAndStream(stderrPipe, "STDERR")

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("执行错误: %v", err)
	}
	return nil
}

// 教师端运行算法并将交互历史推送到 historyChannel
func runAlgorithmScript_tea(exeData ExeData, scriptPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	// 从完整路径提取脚本名
	scriptName := filepath.Base(scriptPath)

	cmd := exec.CommandContext(
		ctx,
		"venv\\Scripts\\python.exe", 
		"-u",
		scriptName,
		"--type_of_prompt", "json_text",
		"--exercise_type", exeData.ExerciseType,
		"--output_type", "json",
		"--Number_of_Generations", exeData.ExerciseQuantity,
		"--Knowledge_Concept", exeData.KnowledgePoint,
		"--Language", exeData.Language,
		"--Difficulty", exeData.Difficulty,
		"--Subject", exeData.Subject,
	)
	cmd.Dir = "../exeGen"

	// 创建管道
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("创建stdout管道失败: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("创建stderr管道失败: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动失败: %v", err)
	}

	// 实时扫描输出
	scanAndStream := func(pipe io.Reader, prefix string) {
		reader := bufio.NewReader(pipe)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					if err != io.EOF {
						log.Printf("%s读取错误: %v", prefix, err)
					}
					return
				}
				msg := strings.TrimSuffix(line, "\n")
				if msg != "" {
					trimmedMsg := strings.Replace(msg, prefix, "", 1)
					historyChannel <- trimmedMsg // 推送处理后的消息
				}
			}
		}
	}

	// 启动输出处理协程
	go scanAndStream(stdoutPipe, "STDOUT")
	go scanAndStream(stderrPipe, "STDERR")

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("执行错误: %v", err)
	}
	return nil
}

// WebSocket处理
func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}
	defer conn.Close() // 确保最终关闭连接

	// 注册客户端
	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()
	log.Printf("新客户端连接，当前总数: %d", len(clients))
	// 心跳配置
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// 发送心跳
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()
	// 保持连接活跃
	for {
		// 持续读取消息（这里只需保持连接）
		if _, _, err := conn.NextReader(); err != nil {
			// 客户端断开时清理
			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()
			break
		}
	}
}

// 启动一个协程，负责从 historyChannel 中读取历史消息并广播给所有客户端
func broadcastHistory() {
	for history := range historyChannel {
		clientsMutex.Lock()
		// 复制当前客户端列表防止并发修改
		currentClients := make([]*websocket.Conn, 0, len(clients))
		for client := range clients {
			currentClients = append(currentClients, client)
		}
		clientsMutex.Unlock()
		// 异步发送避免阻塞
		for _, client := range currentClients {
			go func(c *websocket.Conn) {
				clientsMutex.Lock()
				defer clientsMutex.Unlock()
				if _, ok := clients[c]; !ok {
					return // 连接已关闭
				}
				if err := c.WriteMessage(websocket.TextMessage, []byte(history)); err != nil {
					log.Printf("发送失败: %v", err)
					c.Close()
					delete(clients, c)
				}
			}(client)
		}
	}
}
func main() {
	go broadcastHistory()
	config := lib.LoadConfig[Config]()
	db := lib.NewDB(&config.DatabaseConfig, func(db *gorm.DB) error {
		return db.AutoMigrate(&Graph{}, &Node{}, &Link{}, &DuplicateNode{})
	})

	r := gin.Default()
	r.Use(cors.Default())
	lib.APIBuilder(func(rg *gin.RouterGroup) *gin.RouterGroup {
		rg.GET("", lib.GetAll[Graph](db, nil))
		rg.GET("/by_graph/:gid", lib.GetAll[Graph](db, func(d *gorm.DB, ctx *gin.Context) *gorm.DB {
			return d.Where("graph_id = ?", ctx.Param("gid"))
		}))
		rg.POST("", lib.Create[Graph](db, nil))
		rg.POST("/new_kg", func(ctx *gin.Context) {
			req := new(struct {
				Graph Graph
				Nodes []Node
				Links []Link
			})
			if err := ctx.ShouldBindJSON(req); err != nil {
				log.Println("[new_kg] parse data err: ", err)
				ctx.AbortWithStatus(500)
			} else {
				// Get the current maximum Id values from the database
				var maxNodeId, maxLinkId uint32
				db.Model(&Node{}).Select("MAX(id)").Scan(&maxNodeId)
				db.Model(&Link{}).Select("MAX(id)").Scan(&maxLinkId)

				// Increment Id values for new entries
				newNodeId := maxNodeId + 1
				nodeIdMap := make(map[uint32]uint32) // Map to track old to new Id mapping

				for i := range req.Nodes {
					req.Nodes[i].Id = newNodeId
					nodeIdMap[uint32(i+1)] = newNodeId // Assuming original Ids start from 1
					newNodeId++
				}

				newLinkId := maxLinkId + 1
				for i := range req.Links {
					req.Links[i].Id = newLinkId
					req.Links[i].Source = nodeIdMap[req.Links[i].Source]
					req.Links[i].Target = nodeIdMap[req.Links[i].Target]
					newLinkId++
				}

				db.Create(&req.Graph)
				for i := range req.Nodes {
					req.Nodes[i].GraphId = req.Graph.Id
					db.Create(&req.Nodes[i])
				}
				for i := range req.Links {
					req.Links[i].GraphId = req.Graph.Id
					db.Create(&req.Links[i])
				}
				ctx.JSON(200, req)
			}
		})
		rg.PUT("/:id", lib.Update[Graph](db))
		rg.DELETE("/:id", lib.Delete[Graph](db))
		return rg
	})(r, "/graph")
	lib.AddCRUDNew[Node](r, "/node", db, nil, func(d *gorm.DB, ctx *gin.Context) *gorm.DB {
		return d.Where("graph_id = ?", ctx.Param("gid"))
	}, nil)
	lib.APIBuilder(func(rg *gin.RouterGroup) *gin.RouterGroup {
		rg.GET("", lib.GetAll[Link](db, nil))
		rg.GET("/by_graph/:gid", lib.GetAll[Link](db, func(d *gorm.DB, ctx *gin.Context) *gorm.DB {
			return d.Where("graph_id = ?", ctx.Param("gid"))
		}))
		rg.POST("", lib.Create[Link](db, nil))
		rg.PUT("/:id", lib.Update[Link](db))
		rg.DELETE("/:id", lib.Delete[Link](db))
		return rg
	})(r, "/link")
	lib.APIBuilder(func(rg *gin.RouterGroup) *gin.RouterGroup {
		rg.GET("", lib.GetAll[Link](db, nil))
		rg.GET("/:gid_a/:gid_b", lib.GetAll[Link](db, func(d *gorm.DB, ctx *gin.Context) *gorm.DB {
			return d.Where("graph_id_1 = ? AND graph_id_2 = ?", ctx.Param("gid_a"), ctx.Param("gid_b"))
		}))
		rg.POST("", lib.Create[Link](db, nil))
		rg.PUT("/:id", lib.Update[Link](db))
		rg.DELETE("/:id", lib.Delete[Link](db))
		return rg
	})(r, "/duplicated")
	// 提供 WebSocket 路由
	r.GET("/api/ws/history", handleWebSocket)
	r.GET("/api/ExeAlgorithm-stu", func(c *gin.Context) {
		params := map[string]string{
			"type":   c.Query("exerciseType"),
			"qty":    c.Query("exerciseQuantity"),
			"model":  c.Query("selectedModel"),
			"agents": c.Query("selectedAgents"),
		}
		log.Printf("学生端收到请求 | 参数: %+v", params)

		// 参数校验
		if params["type"] == "" || params["qty"] == "" {
			log.Printf("参数校验失败 | 缺失必要参数")
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "缺少必要参数",
				"data":    nil,
			})
			return
		}

		// 处理请求
		result := handleAlgorithm(ExeData{
			ExerciseType:     params["type"],
			ExerciseQuantity: params["qty"],
			SelectedModel:    params["model"],
			SelectedAgents:   params["agents"],
		})

		// log.Printf("算法处理结果: %s", result)
		// 统一响应格式
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data":    result,
		})
	})
	r.GET("/api/ExeAlgorithm-tea", func(c *gin.Context) {
		params := map[string]string{
			"type":           c.Query("exerciseType"),
			"qty":            c.Query("exerciseQuantity"),
			"model":          c.Query("selectedModel"),
			"agents":         c.Query("selectedAgents"),
			"difficulty":     c.Query("difficulty"),
			"knowledgePoint": c.Query("knowledgePoint"),
			"language":       c.Query("language"),
			"subject":        c.Query("subject"),
		}
		log.Printf("教师端收到请求 | 参数: %+v", params)

		// 参数校验
		if params["type"] == "" || params["qty"] == "" {
			log.Printf("参数校验失败 | 缺失必要参数")
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "缺少必要参数",
				"data":    nil,
			})
			return
		}

		// 处理请求
		result := handleAlgorithm_tea(ExeData{
			ExerciseType:     params["type"],
			ExerciseQuantity: params["qty"],
			SelectedModel:    params["model"],
			SelectedAgents:   params["agents"],
			Difficulty:       params["difficulty"],
			KnowledgePoint:   params["knowledgePoint"],
			Language:         params["language"],
			Subject:          params["subject"],
		})

		// log.Printf("算法处理结果: %s", result)
		// 统一响应格式
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data":    result,
		})
	})
	r.Run(":8010")
}
