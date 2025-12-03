import { baseUrlApi } from "@/router/utils";
import { http } from "@/utils/http";

export type Graph = {
  id: string;
  name: string;
};

export type Node = {
  id: string;
  graph_id: number;
  name: string;
  type: string;
  src: string;
};

export type Link = {
  id: number;
  graph_id: number;
  source: number;
  target: number;
  type: string;
};

/** 获取图谱列表 */
export const apiGetGraphLists = (data?: object) => http.request<Array<Graph>>("get", baseUrlApi("/graph"), { data });
export const apiGetGraphNodes = (graph_id: string) => http.request<Array<Node>>("get", baseUrlApi(`/node/by_graph/${graph_id}`));
export const apiGetGraphLinks = (graph_id: string) => http.request<Array<Link>>("get", baseUrlApi(`/link/by_graph/${graph_id}`));
export const apiGetAllNodes = () => http.request<Array<Node>>("get", baseUrlApi(`/node`));
export const apiGetAllLinks = () => http.request<Array<Link>>("get", baseUrlApi(`/link`));

export const apiDeleteGraph = (graph_id: string) => http.request<Array<Graph>>("delete", baseUrlApi(`/graph/${graph_id}`));
export const apiDeleteNode = (node_id: string) => http.request<Array<Node>>("delete", baseUrlApi(`/node/${node_id}`));
export const apiDeleteLink = (link_id: string) => http.request<Array<Link>>("delete", baseUrlApi(`/link/${link_id}`));

export const apiSetGraph = (graph_id: string, data?: object) => http.request<Array<Graph>>("put", baseUrlApi(`/graph/${graph_id}`), { data });
export const apiSetNode = (node_id: string) => http.request<Array<Node>>("post", baseUrlApi(`/graph/${node_id}/nodes`));
export const apiSetLink = (link_id: string) => http.request<Array<Link>>("post", baseUrlApi(`/graph/${link_id}/links`));
export const apiSetDuplicate = (node_1: string, node_2: string) => http.request<Array<Node>>("post", baseUrlApi(`/repeatNode/${node_1}/${node_2}`));

export const apiAddGraph = (data?: object) => http.request<Array<Graph>>("post", baseUrlApi("/graphs"), { data });
export const apiAddNode = (graph_id: string, data?: object) => http.request<Array<Node>>("post", baseUrlApi(`/graph/${graph_id}/nodes`), { data });
export const apiAddLink = (graph_id: string) => http.request<Array<Link>>("post", baseUrlApi(`/graph/${graph_id}/links`));
export const apiAddDuplicate = (graph_1: string, graph_2: string) => http.request<Array<Node>>("post", baseUrlApi(`/duplicates/${graph_1}/${graph_2}`));

export const apiNewKG = (data?: object) => http.request<Array<Graph>>("post", baseUrlApi("/graph/new_kg"), { data });