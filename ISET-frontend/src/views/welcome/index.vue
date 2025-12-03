<template>
  <div class="overall">
    <div class="card-wrap" @mousemove="handleMouseMove" @mouseenter="handleMouseEnter" @mouseleave="handleMouseLeave"
      ref="card">
      <div class="card" :style="cardStyle">
        <img
          class="card-bg"
          src="@/assets/svg/welcome2.jpg"
          :style="cardBgTransform"
          alt
        />
      </div>
    </div>
  </div>
</template>

<script>
export default {
  computed: {
    mousePX() {
      return this.mouseX / this.width;
    },
    mousePY() {
      return this.mouseY / this.height;
    },
    cardStyle() {
      // 外层卡片偏移量
      const rX = this.mousePX * 30;
      const rY = this.mousePY * -30;
      return {
        transform: `rotateY(${rX}deg) rotateX(${rY}deg)`
      };
    },
    cardBgTransform() {
      // 内层图片偏移量
      const tX = this.mousePX * -40;
      const tY = this.mousePY * -40;
      return {
        transform: `translateX(${tX}px) translateY(${tY}px)`
      };
    }
  },
  data() {
    return {
      width: 0,
      height: 0,
      mouseX: 0,
      mouseY: 0,
      mouseLeaveDelay: null //添加一个默认还原的定时器
    };
  },
  mounted() {
    this.width = this.$refs.card.offsetWidth;
    this.height = this.$refs.card.offsetHeight;
  },
  methods: {
    handleMouseMove(e) {
      // e.pageX 鼠标在body的x轴坐标位置
      //this.$refs.card.offsetLeft 外层卡片容器与body的左侧距离
      // e.pageX - this.$refs.card.offsetLeft 鼠标所在容器内的x轴坐标位置
      // e.pageX - this.$refs.card.offsetLeft - this.width / 2 鼠标以容器中心点为原点的x轴位置
      this.mouseX = e.pageX - this.$refs.card.offsetLeft - this.width / 2;
      this.mouseY = e.pageY - this.$refs.card.offsetTop - this.height / 2;
    },
    //当鼠标进入时清除计时器
    handleMouseEnter() {
      clearTimeout(this.mouseLeaveDelay);
    },
    //当鼠标离开时添加定时器
    handleMouseLeave() {
      this.mouseLeaveDelay = setTimeout(() => {
        this.mouseX = 0;
        this.mouseY = 0;
      }, 1000);
    }
  }
};
</script>

<style lang="scss" scoped>
.overall {
  perspective: 2000;
}

.card-wrap {
  perspective: 2000;
  display: flex;
  transform-style: preserve-3d;
  padding: 50px 100px 0px 100px;
  cursor: pointer;
}

.card {
  transform-style: preserve-3d;
  width: 100%;
  height: 100%;
  // width: 1500px;
  // height: 700px;
  flex: 1;
  background: #ececec;
  overflow: hidden;
  perspective: 1000;
  transition: 1s cubic-bezier(0.23, 1, 0.32, 1);
  opacity: 0.8;
  object-fit: contain;
}

.card-bg {
  transform-style: preserve-3d;
  width: 290%;
  height: 135%;
  display: block;
  justify-content: center;
  align-items: center;
  // position: relative;
  // left: -0%;
  // top: -25%;
  padding: auto;
  transition: 1s cubic-bezier(0.23, 1, 0.32, 1);
  //   pointer-events: none; //pointer-events 更像是JavaScript。用于阻止鼠标事件
}

//样式尽量定义成3层（否则会出现卡顿现象），1.鼠标经过开始层，2.卡片旋转层，3.图片旋转层
//鼠标移入到1层时，执行过度效果（会有一个缓冲），过1s后缓冲结束，此时旋转时不再会有以上定义的1s过渡效果（不然会出现一卡一卡的情况），当鼠标移出时才会执行以上的过度效果，实现缓慢的归位状态（而不是一闪就归位了）
//此处 hover选择器优先级最高
.card-wrap:hover .card-bg {
  transition: 0.5s cubic-bezier(0.23, 1, 0.32, 1);
}

.card-wrap:hover .card {
  opacity: 1;
  transition: 0.5s cubic-bezier(0.23, 1, 0.32, 1);
  box-shadow: white 0 0 0 5px, rgba(0, 0, 0, 0.66) 0 30px 60px 0;
}
</style>
