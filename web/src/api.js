/*global process */
const Endpoint = process.env.NODE_ENV === 'production' ?
  'https://xxxxxxxxxx.execute-api.ap-south-1.amazonaws.com/v1/graphql' :
  'http://localhost:8080/graphql';

const request = async (operation, body) => {
  const response = await fetch(Endpoint, {
    method: 'POST',
    mode: 'cors',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json;charset=utf-8'
    },
    body: JSON.stringify(body)
  });

  const { data: { [operation]: result } } = await response.json();

  return result;
};

const Api = {
  list: async (model) => {
    const variables = {};

    if (model) {
      if (model.status) {
        variables.status = model.status;
      }

      if (model.dueAt) {
        variables.dueAt = model.dueAt;
      }

      if (model.startKey) {
        variables.startKey = model.startKey;
      }
    }

    const body = {
      query: `
        query List($status: ScheduleStatus, $dueAt: DateRange, $startKey: ScheduleListStartKey) {
          list(status: $status, dueAt: $dueAt, startKey: $startKey) {
            schedules {
              id
              dueAt
              url
              method
              status
            }
            nextKey {
              id,
              dueAt,
              status
            }
          }
        }
      `,
      variables
    };

    return request('list', body);
  },

  get: async (id) => {
    const body = {
      query: `
        query Get($id: ID!) {
          get(id: $id) {
            id
            dueAt
            url
            method
            headers
            body
            status
            startedAt
            completedAt
            canceledAt
            result
            createdAt
          }
        }
      `,
      variables: { id }
    };

    return request('get', body);
  },

  create: async (model) => {
    const body = {
      query: `
        mutation Create($dueAt: DateTime!, $url: String!, $method: HTTPMethod!, $headers: StringMap, $body: String) {
          create(dueAt: $dueAt, url: $url, method: $method, headers: $headers, body: $body)
        }
      `,
      variables: model
    };

    return request('create', body);
  },

  cancel: async (id) => {
    const body = {
      query: `
        mutation Cancel($id: ID!) {
          cancel(id: $id)
        }
      `,
      variables: { id }
    };

    return request('cancel', body);
  }
};

export default Api;
