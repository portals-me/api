const sdk = (endpoint, axios, firebase) => {
  return {
    project: {
      list: async () => {
        return (await axios.get(`${endpoint}/projects`)).data;
      },
      get: async (projectId) => {
        return (await axios.get(`${endpoint}/projects/${projectId}`)).data;
      },
    },
    user: {
      getMe: async () => {
        return (await axios.get(`${endpoint}/users/me`)).data;
      },
    },
  };
};

export default sdk;
