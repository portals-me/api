const genSDK = (url, token, axios) => ({
  signIn: async (gtoken) => {
    const result = await axios.post(`${url}/signIn`, gtoken);
    return result.data;
  },
  user: {
    me: async () => {
      const result = await axios.get(`${url}/users/me`, { headers: { Authorization: `Bearer ${token}` } });
      return result.data;
    },
  },
  comment: {
    create: async (projectId, message) => {
      const result = await axios.post(`${url}/comments`, {
        projectId,
        message,
      }, { headers: { Authorization: `Bearer ${token}` } });
      return result.data;
    },
    list: async (projectId) => {
      const result = await axios.post(`${url}/projects/${projectId}/comments`, {
        projectId,
      }, { headers: { Authorization: `Bearer ${token}` } });
      return result.data;
    },
  },
  project: {
    create: async (form) => {
      return await axios.post(
        `${url}/projects`,
        form,
        { headers: { Authorization: `Bearer ${token}` } }
      );
    },
    get: async (projectId) => {
      return await axios.get(`${url}/projects/${projectId}`, { headers: { Authorization: `Bearer ${token}` } });
    },
    list: async () => {
      return await axios.get(`${url}/projects`, { headers: { Authorization: `Bearer ${token}` } });
    },
  },
});

module.exports = {
  genSDK,
};
