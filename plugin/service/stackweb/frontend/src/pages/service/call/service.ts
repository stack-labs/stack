import request from 'umi-request';
import {Pagination, CallParam} from './data.d';

export async function queryServices(params: Pagination) {
  const data = request('/platform/api/v1/b/services', {
    method: 'POST',
    data: {
      ...params,
    },
  });

  return data;
}

export async function callService(params: CallParam) {
  const data = request('/platform/api/v1/b/rpc', {
    method: 'POST',
    data: {
      ...params,
    },
  });

  return data;
}
