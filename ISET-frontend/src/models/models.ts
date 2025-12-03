// 数据分为两大类
// 数据实体/restful接口请求返回数据

// 数据实体定义
interface Student {
  id: number;
}

interface Course {
  id: number;
  name: string;
  teacher: Teacher;
  students: Student[];
  notices: Notice[];
  chapters: Chapter[];
}

interface Score {
  id: number;
}

interface Teacher {
  id: number;
}

interface Notice {
  id: number;
}

interface Note {
  id: number;
}

interface Chapter {
  attachment: Attachment[];
  content: string;
  id: string;
  title: string;
  [property: string]: any;
}

interface Attachment {
  id: string;
  name: string;
  type: string;
  url: string;
  [property: string]: any;
}
