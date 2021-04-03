const request = async (operation, body) => {
  const response = await fetch('https://nokxy6r6uf.execute-api.ap-south-1.amazonaws.com/v1/graphql', {
    method: 'POST',
    mode: 'cors',
    headers: {
      'Content-Type': 'application/json;charset=utf-8'
    },
    body: JSON.stringify(body)
  });

  const { [operation]: data } = await response.json();

  return data;
};

const Api = {
  list: async (model) => {
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
        }`,
      variables: {}
    };

    if (model) {
      if (model.status) {
        body.variables.status = model.status;
      }

      if (model.dueAt) {
        body.variables.dueAt = model.dueAt;
      }

      if (model.startKey) {
        body.variables.startKey = model.startKey;
      }
    }

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
      variables: {
        id
      }
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
      variables: {
        id
      }
    };

    return request('cancel', body);
  }
};

export default Api;
