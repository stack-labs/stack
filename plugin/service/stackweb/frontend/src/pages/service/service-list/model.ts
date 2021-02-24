import {Effect, Reducer} from 'umi';
import {Node, Service} from './data.d';
import {queryServices} from './service';

export interface ServicesState {
  list: Service[];
  filters: FiltersState;
}

export interface FiltersState {
  service: string;
  node: string;
}

export interface ModelType {
  namespace: string;
  state: ServicesState;
  effects: {
    fetch: Effect;
  };
  reducers: {
    queryList: Reducer<ServicesState>;
  };
}

const filterNode = (filter: string, nodes: Node[]) => {
  const nodesTemp: Node[] = [];
  if (filter != null) {
    nodes.forEach((n: Node) => {
      if (
        n.id.indexOf(filter) > 0 ||
        n.address.indexOf(filter) > 0 ||
        JSON.stringify(n.metadata).indexOf(filter) > 0
      ) {
        nodesTemp.push(n);
      }
    });

    return nodesTemp;
  }

  return nodes;
};

const filterService = (filter: string, services: Service[]) => {
  const servicesTemp: Service[] = [];
  if (filter != null && filter !== '') {
    services.forEach((item: Service) => {
      if (item.name.indexOf(filter) > 0) {
        servicesTemp.push(item);
      }
    });

    return servicesTemp;
  }

  return services;
};

const emptyArray = (arr: any[]) => {
  while (arr.length > 0) {
    arr.pop();
  }
};

// @ts-ignore
// @ts-ignore
const Model: ModelType = {
  namespace: 'searchServices',

  state: {
    list: [],
    filters: {
      service: "",
      node: "",
    },
  },

  effects: {
    * fetch({payload}, {call, put}) {
      const response = yield call(queryServices, payload);
      const data = Array.isArray(response.data) ? response.data : [];
      // filter locally
      const {serviceStr, nodeStr} = payload;
      const services: Service[] = filterService(serviceStr, data);

      if (nodeStr != null && nodeStr !== '') {
        services.forEach((service: Service) => {
          emptyArray(service.nodes);
          service.nodes.push(...filterNode(nodeStr, service.nodes));
        });
      }

      yield put({
        type: 'queryList',
        payload: services,
      });
    }
  },

  reducers: {
    queryList(state, action) {
      return {
        ...(state as ServicesState),
        list: action.payload,
      };
    },
  },
};

export default Model;
