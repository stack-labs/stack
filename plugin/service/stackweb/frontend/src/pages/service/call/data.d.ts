export interface Service {
  name: string;
  version: string;
  metadata: Map;
  endpoints: Endpoint[];
  nodes: Node[];
}

export interface Endpoint {
}

export interface Node {
  id: string;
  address: string;
  metadata: Map;
}

export interface Pagination {
  service: Service;
  total: number;
  pageSize: number;
  current: number;
}

export interface PageData {
  data: Service[];
  pagination: Partial<Pagination>;
}

export interface CallParam {
  service: string;
  address: string;
  endpoint: string;
}
