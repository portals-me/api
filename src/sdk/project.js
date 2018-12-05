export const api = (endpoint, axios) => {
  return {
    list: async () => {
      return (await axios.get(`${endpoint}/projects`)).data;
    },
    get: async (projectId) => {
      return (await axios.get(`${endpoint}/projects/${projectId}`)).data;
    },
  };
};
