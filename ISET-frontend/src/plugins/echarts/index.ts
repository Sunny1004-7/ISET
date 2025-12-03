// import type { App } from "vue";
// import * as echarts from "echarts/core";
// import { CanvasRenderer } from "echarts/renderers";
// import { PieChart, BarChart, LineChart } from "echarts/charts";
// import {
//   GridComponent,
//   TitleComponent,
//   LegendComponent,
//   GraphicComponent,
//   ToolboxComponent,
//   TooltipComponent,
//   DataZoomComponent,
//   VisualMapComponent
// } from "echarts/components";

// const { use } = echarts;

// use([
//   PieChart,
//   BarChart,
//   LineChart,
//   CanvasRenderer,
//   GridComponent,
//   TitleComponent,
//   LegendComponent,
//   GraphicComponent,
//   ToolboxComponent,
//   TooltipComponent,
//   DataZoomComponent,
//   VisualMapComponent
// ]);

// /**
//  * @description 按需引入echarts
//  * @see {@link https://echarts.apache.org/handbook/zh/basics/import#%E6%8C%89%E9%9C%80%E5%BC%95%E5%85%A5-echarts-%E5%9B%BE%E8%A1%A8%E5%92%8C%E7%BB%84%E4%BB%B6}
//  * @see 温馨提示：必须将 `$echarts` 添加到全局 `globalProperties` ，为了配合 https://pure-admin-utils.netlify.app/hooks/useEcharts/useEcharts.html 使用
//  */
// export function useEcharts(app: App) {
//   app.config.globalProperties.$echarts = echarts;
// }

// export default echarts;
// 引入 echarts 核心模块。
import * as echarts from "echarts/core";
//引入柱状图和折线图组件。
import { BarChart, LineChart, PieChart, GraphChart } from "echarts/charts";
// 引入标题、提示框、网格、数据集和数据转换器组件。
import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  // 数据集组件
  DatasetComponent,
  // 内置数据转换器组件 (filter, sort)
  TransformComponent
} from "echarts/components";
//引入标签布局和通用过渡动画特性。
import { LabelLayout, UniversalTransition } from "echarts/features";
// 引入 Canvas 渲染器。
import { CanvasRenderer } from "echarts/renderers";

import type {
  // 系列类型的定义后缀都为 SeriesOption
  BarSeriesOption,
  LineSeriesOption
} from "echarts/charts";

import type {
  // 组件类型的定义后缀都为 ComponentOption
  TitleComponentOption,
  TooltipComponentOption,
  GridComponentOption,
  DatasetComponentOption
} from "echarts/components";
import type { ComposeOption } from "echarts/core";

// 通过 ComposeOption 来组合出一个只有必须组件和图表的 Option 类型
type ECOption = ComposeOption<
  | BarSeriesOption
  | LineSeriesOption
  | TitleComponentOption
  | TooltipComponentOption
  | GridComponentOption
  | DatasetComponentOption
>;

/** 
    注册必须的组件，包括标题、提示框、网格、数据集、数据转换器，
    以及柱状图、折线图、标签布局、通用过渡动画和 Canvas 渲染器。
*/
echarts.use([
  TitleComponent,
  TooltipComponent,
  GridComponent,
  DatasetComponent,
  TransformComponent,
  BarChart,
  LineChart,
  PieChart,
  GraphChart,
  LabelLayout,
  UniversalTransition,
  CanvasRenderer,
  LegendComponent
]);
// 导出
export default echarts;
