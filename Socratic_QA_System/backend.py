from flask import Flask, request, jsonify
from flask_cors import CORS
import threading
import time
import uuid
from datetime import datetime
from typing import Dict, Any, List, Optional
import json
import traceback

# 导入你原有的核心组件
from core import EventBus, MessageType, LLMManager, Logger, ConversationOrchestrator, SampleDataManager, Message
from teacher_agent import TeacherAgent
from monitor_agent import MonitorAgent
from knowledge_state_agent import KnowledgeStateAgent
from config import Config

app = Flask(__name__)
CORS(app)  # 允许跨域请求

class SocraticBackend:
    """苏格拉底教学后端服务"""
    
    def __init__(self):
        print(" 开始初始化苏格拉底后端服务...")
        
        # 加载配置
        try:
            Config.load_from_env()
            self.system_config = Config.get_system_config()
            print("✅ 配置加载成功")
        except Exception as e:
            print(f"❌ 配置加载失败: {e}")
            raise
        
        # 使用统一的LLM管理器和日志记录器
        try:
            self.llm_manager = LLMManager(Config.get_llm_config())
            self.logger = Logger(enable_file_logging=False)
            self.event_bus = EventBus(self.logger)
            print(" 核心组件初始化成功")
        except Exception as e:
            print(f" 核心组件初始化失败: {e}")
            raise
        
        # 初始化对话协调器
        try:
            self.conversation_orchestrator = ConversationOrchestrator(self.event_bus, self.logger)
            print(" 对话协调器初始化成功")
        except Exception as e:
            print(f"对话协调器初始化失败: {e}")
            raise
        
        # 初始化智能体
        try:
            self.teacher_agent = TeacherAgent("teacher", self.llm_manager, self.logger)
            print(" 教师智能体创建成功")
            
            # 根据配置决定是否启用监控和知识分析智能体
            if self.system_config.get("enable_monitoring", True):
                self.monitor_agent = MonitorAgent("monitor", self.llm_manager, self.logger)
                print(" 监控智能体创建成功")
            else:
                self.monitor_agent = None
                print(" 监控智能体已禁用")
                
            if self.system_config.get("enable_knowledge_analysis", True):
                self.knowledge_state_agent = KnowledgeStateAgent("knowledge_state", self.llm_manager, self.logger)
                print(" 知识状态智能体创建成功")
            else:
                self.knowledge_state_agent = None
                print(" 知识状态智能体已禁用")
        except Exception as e:
            print(f" 智能体初始化失败: {e}")
            raise
        
        # 注册智能体到事件总线
        try:
            self.event_bus.register_agent(self.teacher_agent)
            if self.monitor_agent:
                self.event_bus.register_agent(self.monitor_agent)
            if self.knowledge_state_agent:
                self.event_bus.register_agent(self.knowledge_state_agent)
            print("智能体注册到事件总线成功")
        except Exception as e:
            print(f" 智能体注册失败: {e}")
            raise
        
        # 启动事件总线
        try:
            self.event_bus.start()
            print(" 事件总线启动成功")
        except Exception as e:
            print(f" 事件总线启动失败: {e}")
            raise
        
        # 初始化所有智能体
        try:
            self.teacher_agent.initialize()
            if self.monitor_agent:
                self.monitor_agent.initialize()
            if self.knowledge_state_agent:
                self.knowledge_state_agent.initialize()
            print(" 智能体初始化完成")
        except Exception as e:
            print(f"智能体初始化失败: {e}")
            raise
        
        # 启动所有智能体
        try:
            self.teacher_agent.start()
            if self.monitor_agent:
                self.monitor_agent.start()
            if self.knowledge_state_agent:
                self.knowledge_state_agent.start()
            print(" 智能体启动完成")
        except Exception as e:
            print(f" 智能体启动失败: {e}")
            raise
        
        # 等待知识状态智能体初始化完成
        if self.knowledge_state_agent:
            self._wait_for_knowledge_agent_ready()
        
        # 会话管理
        self.conversations = {}  # 存储多个会话
        
        # 知识状态摘要
        try:
            self.knowledge_summary = None
            self.exercise_records = SampleDataManager.get_sample_exercise_records()
            self._get_knowledge_summary()
            print(" 知识状态摘要获取完成")
        except Exception as e:
            print(f" 知识状态摘要获取失败: {e}")
            self.knowledge_summary = "知识状态分析初始化失败"
        
        # 分析默认记录
        if self.knowledge_state_agent:
            self._analyze_default_exercise_records()
        
        print(" 苏格拉底后端服务初始化完成!")
    
    def _wait_for_knowledge_agent_ready(self):
        """等待知识状态智能体准备就绪"""
        max_wait_time = 30
        wait_interval = 0.5
        elapsed_time = 0
        
        print("等待知识状态智能体初始化...")
        self.logger.log_system("等待知识状态智能体初始化...")
        
        while elapsed_time < max_wait_time:
            if (self.knowledge_state_agent and 
                hasattr(self.knowledge_state_agent, 'analysis_ready') and
                self.knowledge_state_agent.analysis_ready):
                print(" 知识状态智能体已就绪")
                self.logger.log_system("知识状态智能体已就绪")
                return
            time.sleep(wait_interval)
            elapsed_time += wait_interval
        
        print(" 知识状态智能体初始化超时，继续启动系统")
        self.logger.log_system("知识状态智能体初始化超时，继续启动系统")
    
    def _get_knowledge_summary(self):
        """获取知识状态摘要"""
        if not self.knowledge_state_agent:
            self.knowledge_summary = "知识状态分析功能已禁用"
            return
        
        try:
            if not hasattr(self.knowledge_state_agent, 'analysis_ready') or not self.knowledge_state_agent.analysis_ready:
                self.knowledge_summary = "系统正在初始化知识状态分析，使用默认设置"
                return
            
            if hasattr(self.knowledge_state_agent, 'get_knowledge_summary'):
                self.knowledge_summary = self.knowledge_state_agent.get_knowledge_summary()
            else:
                self.knowledge_summary = "知识状态分析功能可用"
            
            if not self.knowledge_summary:
                self.knowledge_summary = "知识状态分析完成，系统已就绪"
                
        except Exception as e:
            print(f" 获取知识状态摘要失败: {e}")
            self.logger.log_error("SYSTEM", f"获取知识状态摘要失败: {e}")
            self.knowledge_summary = "知识状态分析初始化失败，使用默认设置"
    
    def _analyze_default_exercise_records(self):
        """调用知识状态追踪智能体分析默认记录"""
        if not self.knowledge_state_agent or not self.exercise_records:
            print(" 知识状态智能体或习题记录未准备好，跳过分析")
            return
        
        try:
            exercise_data_json = json.dumps(self.exercise_records)
            
            message = Message(
                id=str(uuid.uuid4()),
                sender="system",
                recipient="knowledge_state",
                type=MessageType.DATA_REQUEST,
                content={"exercise_data": exercise_data_json},
                timestamp=datetime.now().isoformat()
            )
            self.event_bus.send_message(message)
            print("✅ 习题记录已发送给知识状态智能体分析")
            
        except Exception as e:
            print(f" 发送习题记录失败: {e}")
            self.logger.log_error("SYSTEM", f"发送习题记录失败: {e}")
    
    def create_conversation(self) -> str:
        """创建新会话"""
        conversation_id = str(uuid.uuid4())
        self.conversations[conversation_id] = {
            "id": conversation_id,
            "round_number": 0,
            "history": [],
            "created_at": datetime.now().isoformat(),
            "mode": "fast"  # 默认快思考模式
        }
        print(f" 创建新会话: {conversation_id}")
        return conversation_id
    
    def slow_think_response(self, conversation_id: str, user_message: str) -> str:
        """慢思考模式 - 使用复杂智能体系统"""
        print(f"开始慢思考处理 - 会话: {conversation_id}")
        print(f"用户消息: {user_message}")
        
        if conversation_id not in self.conversations:
            print(f" 会话不存在: {conversation_id}")
            return "会话不存在"
        
        conversation = self.conversations[conversation_id]
        conversation["round_number"] += 1
        
        # 记录用户消息
        conversation["history"].append({
            "role": "user",
            "content": user_message,
            "round": conversation["round_number"],
            "timestamp": datetime.now().isoformat()
        })
        
        print(f"记录用户消息，当前轮次: {conversation['round_number']}")
        
        # 发送消息给教师智能体
        try:
            message = Message(
                id=str(uuid.uuid4()),
                sender="system",
                recipient="teacher",
                type=MessageType.TASK_REQUEST,
                content={
                    "conversation_id": conversation_id,
                    "student_message": user_message,
                    "round_number": conversation["round_number"],
                    "student_state": {
                        "conversation_history": conversation["history"],
                        "knowledge_summary": self.knowledge_summary
                    }
                },
                timestamp=datetime.now().isoformat()
            )
            self.event_bus.send_message(message)
            print(" 消息已发送给教师智能体")
        except Exception as e:
            print(f" 发送消息给教师智能体失败: {e}")
            return f"发送消息失败: {e}"
        
        # 等待教师回复
        print(" 等待教师智能体回复...")
        teacher_response = self._wait_for_teacher_response()
        
        if teacher_response:
            print(" 获得教师回复")
            # 记录教师回复
            conversation["history"].append({
                "role": "teacher",
                "content": teacher_response,
                "round": conversation["round_number"],
                "timestamp": datetime.now().isoformat()
            })
            return teacher_response
        else:
            print(" 教师回复超时或失败")
            return "抱歉，我在生成回复时遇到了一些问题。请重新描述一下你的问题。"
    
    def _wait_for_teacher_response(self) -> str:
        """等待教师智能体回复"""
        wait_interval = 0.5
        max_wait_time = 60  # 最大等待60秒
        elapsed_time = 0
        
        while elapsed_time < max_wait_time:
            try:
                # 检查TeacherAgent是否有回复
                if hasattr(self.teacher_agent, 'conversation_history') and self.teacher_agent.conversation_history:
                    # 获取最新的教师回复
                    for message in reversed(self.teacher_agent.conversation_history):
                        if message.get("role") == "teacher":
                            teacher_response = message.get("content", "")
                            if teacher_response and teacher_response != "正在生成回复...":
                                print(f" 教师回复: {teacher_response[:100]}...")
                                return teacher_response
                
                # 显示等待进度
                if elapsed_time % 10 == 0 and elapsed_time > 0:
                    print(f"⏳ 仍在等待教师回复... ({elapsed_time}秒)")
                
                time.sleep(wait_interval)
                elapsed_time += wait_interval
                
            except Exception as e:
                print(f" 检查教师回复时出错: {e}")
                break
        
        print(f" 教师回复等待超时 ({max_wait_time}秒)")
        return None

# 创建全局实例
print(" 正在创建苏格拉底后端实例...")
try:
    socratic_backend = SocraticBackend()
    print(" 后端实例创建成功")
except Exception as e:
    print(f" 后端实例创建失败: {e}")
    print(traceback.format_exc())
    raise

# ==================== Flask API 路由 ====================

@app.route('/api/conversation', methods=['POST'])
def create_conversation():
    """创建新会话"""
    try:
        print(" 收到创建会话请求")
        conversation_id = socratic_backend.create_conversation()
        response_data = {
            "success": True,
            "conversation_id": conversation_id,
            "knowledge_summary": socratic_backend.knowledge_summary
        }
        print(f" 会话创建成功，返回数据: {response_data}")
        return jsonify(response_data)
    except Exception as e:
        error_msg = f"创建会话失败: {e}"
        print(f" {error_msg}")
        print(traceback.format_exc())
        return jsonify({"success": False, "error": error_msg}), 500

@app.route('/api/chat', methods=['POST'])
def chat():
    """处理聊天消息"""
    try:
        print(" 收到聊天请求")
        
        # 验证请求格式
        if not request.is_json:
            error_msg = "请求必须是JSON格式"
            print(f" {error_msg}")
            return jsonify({"success": False, "error": error_msg}), 400
        
        data = request.json
        print(f" 请求数据: {data}")
        
        # 验证必要参数
        conversation_id = data.get('conversation_id')
        message = data.get('message')
        mode = data.get('mode', 'fast')  # 默认快思考
        
        print(f" 参数验证 - 会话ID: {conversation_id}, 消息: {message}, 模式: {mode}")
        
        if not conversation_id:
            error_msg = "缺少conversation_id参数"
            print(f" {error_msg}")
            return jsonify({"success": False, "error": error_msg}), 400
            
        if not message:
            error_msg = "缺少message参数"
            print(f" {error_msg}")
            return jsonify({"success": False, "error": error_msg}), 400
        
        # 验证会话是否存在
        if conversation_id not in socratic_backend.conversations:
            error_msg = f"会话不存在: {conversation_id}"
            print(f" {error_msg}")
            return jsonify({"success": False, "error": error_msg}), 404
        
        if mode == 'slow':
            print(" 进入慢思考模式")
            # 慢思考模式 - 使用复杂智能体系统
            response = socratic_backend.slow_think_response(conversation_id, message)
            response_data = {
                "success": True,
                "response": response,
                "mode": "slow",
                "processing_info": "使用了完整的ICECoT智能体分析流程：情绪分析→意图推断→策略选择→响应生成"
            }
            print(f" 慢思考响应: {response_data}")
            return jsonify(response_data)
        else:
            print(" 快思考模式 - 返回标识给前端处理")
            # 快思考模式 - 返回标识让前端使用原有prompt方式
            return jsonify({
                "success": True,
                "mode": "fast",
                "message": "请使用前端快思考模式处理"
            })
            
    except Exception as e:
        error_msg = f"聊天处理失败: {e}"
        print(f" {error_msg}")
        print(traceback.format_exc())
        return jsonify({"success": False, "error": error_msg}), 500

@app.route('/api/conversation/<conversation_id>', methods=['GET'])
def get_conversation(conversation_id):
    """获取会话历史"""
    try:
        print(f" 获取会话历史: {conversation_id}")
        if conversation_id not in socratic_backend.conversations:
            error_msg = "会话不存在"
            print(f" {error_msg}")
            return jsonify({"success": False, "error": error_msg}), 404
        
        conversation = socratic_backend.conversations[conversation_id]
        print(f" 返回会话历史，消息数: {len(conversation.get('history', []))}")
        return jsonify({
            "success": True,
            "conversation": conversation
        })
    except Exception as e:
        error_msg = f"获取会话历史失败: {e}"
        print(f" {error_msg}")
        return jsonify({"success": False, "error": error_msg}), 500

@app.route('/api/conversation/<conversation_id>/mode', methods=['PUT'])
def update_conversation_mode(conversation_id):
    """更新会话模式"""
    try:
        print(f" 更新会话模式: {conversation_id}")
        data = request.json
        mode = data.get('mode', 'fast')
        
        if conversation_id not in socratic_backend.conversations:
            error_msg = "会话不存在"
            print(f" {error_msg}")
            return jsonify({"success": False, "error": error_msg}), 404
        
        socratic_backend.conversations[conversation_id]["mode"] = mode
        print(f"模式更新成功: {mode}")
        return jsonify({
            "success": True,
            "mode": mode
        })
    except Exception as e:
        error_msg = f"更新会话模式失败: {e}"
        print(f" {error_msg}")
        return jsonify({"success": False, "error": error_msg}), 500

@app.route('/api/health', methods=['GET'])
def health_check():
    """健康检查"""
    try:
        print(" 健康检查请求")
        health_data = {
            "status": "healthy",
            "timestamp": datetime.now().isoformat(),
            "agents_ready": {
                "teacher": socratic_backend.teacher_agent is not None,
                "monitor": socratic_backend.monitor_agent is not None,
                "knowledge_state": socratic_backend.knowledge_state_agent is not None
            },
            "knowledge_summary": socratic_backend.knowledge_summary,
            "active_conversations": len(socratic_backend.conversations),
            "system_info": {
                "event_bus_active": socratic_backend.event_bus is not None,
                "llm_manager_ready": socratic_backend.llm_manager is not None
            }
        }
        print(f"健康检查响应: {health_data}")
        return jsonify(health_data)
    except Exception as e:
        error_msg = f"健康检查失败: {e}"
        print(f" {error_msg}")
        return jsonify({"status": "unhealthy", "error": error_msg}), 500

# 错误处理
@app.errorhandler(404)
def not_found(error):
    print(f" 404错误: {request.url}")
    return jsonify({"success": False, "error": "API端点不存在"}), 404

@app.errorhandler(405)
def method_not_allowed(error):
    print(f" 405错误: {request.method} {request.url}")
    return jsonify({"success": False, "error": "HTTP方法不允许"}), 405

@app.errorhandler(500)
def internal_error(error):
    print(f" 500错误: {error}")
    return jsonify({"success": False, "error": "服务器内部错误"}), 500

# 主程序入口
if __name__ == '__main__':
    print("=" * 60)
    print(" 智能苏格拉底教学范式问答系统 - 后端服务")
    print(" 运行地址: http://localhost:5000")
    print(" 调试模式: 启用")
    print("=" * 60)
    
    # 显示系统状态
    try:
        print("\n 系统组件状态:")
        print(f"  教师智能体: {' 就绪' if socratic_backend.teacher_agent else ' 未就绪'}")
        print(f"  监控智能体: {' 就绪' if socratic_backend.monitor_agent else ' 禁用'}")
        print(f"  知识状态智能体: {' 就绪' if socratic_backend.knowledge_state_agent else ' 禁用'}")
        print(f"  知识摘要: {socratic_backend.knowledge_summary}")
        print("")
    except Exception as e:
        print(f" 无法显示系统状态: {e}")
    
    app.run(debug=True, host='0.0.0.0', port=5000)