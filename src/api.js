const genSDK = (url, token, axios) => ({
  signUp: async (json) => {
    return await axios.post(`${url}/signUp`, JSON.stringify(json));
  },
  signIn: async (gtoken) => {
    return await axios.post(`${url}/signIn`, gtoken);
  },
  user: {
    me: async () => {
      return await axios.get(`${url}/users/me`, { headers: { Authorization: `Bearer ${token}` } });
    },
  },
  comment: {
    create: async (collectionId, message) => {
      return await axios.post(`${url}/comments`, {
        collectionId,
        message,
      }, { headers: { Authorization: `Bearer ${token}` } });
    },
    list: async (collectionId) => {
      return await axios.get(`${url}/collections/${collectionId}/comments`, { headers: { Authorization: `Bearer ${token}` } });
    },
  },
  collection: {
    create: async (form) => {
      return await axios.post(
        `${url}/collections`,
        form,
        { headers: { Authorization: `Bearer ${token}` } }
      );
    },
    get: async (collectionId) => {
      return await axios.get(`${url}/collections/${collectionId}`, { headers: { Authorization: `Bearer ${token}` } });
    },
    list: async () => {
      return await axios.get(`${url}/collections`, { headers: { Authorization: `Bearer ${token}` } });
    },
    delete: async (collectionId) => {
      return await axios.delete(`${url}/collections/${collectionId}`, { headers: { Authorization: `Bearer ${token}` } });
    },
  },
});

module.exports = {
  genSDK,
};
