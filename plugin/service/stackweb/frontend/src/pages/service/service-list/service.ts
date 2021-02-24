import request from 'umi-request';
import { Pagination } from './data.d';

export async function queryServices(params: Pagination) {
  const data = request('/platform/api/v1/b/services', {
    method: 'POST',
    data: {
      ...params,
    },
  });

  return data;
}
