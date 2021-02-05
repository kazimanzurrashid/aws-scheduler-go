const http = async (operation, body) => {
  const res = await fetch('https://api.schedules.my-domain.com/v1/graphql', {
    method: 'POST',
    mode: 'cors',
    headers: {
      'Content-Type': 'application/json;utf-8'
    },
    body: JSON.stringify(body)
  });

  const { [operation]: data } = await res.json();

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

    return http('list', body);
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

    return http('get', body);
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

    return http('create', body);
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

    return http('cancel', body);
  }
};

export default Api;
