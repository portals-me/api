const genSDK = (url, token, axios) => ({
  signIn: async (gtoken) => {
    return await axios.post(`${url}/signIn`, gtoken);
  },
  user: {
    me: async () => {
      return await axios.get(`${url}/users/me`, { headers: { Authorization: `Bearer ${token}` } });
    },
  },
  comment: {
    create: async (projectId, message) => {
      return await axios.post(`${url}/comments`, {
        projectId,
        message,
      }, { headers: { Authorization: `Bearer ${token}` } });
    },
    list: async (projectId) => {
      return await axios.get(`${url}/projects/${projectId}/comments`, { headers: { Authorization: `Bearer ${token}` } });
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
