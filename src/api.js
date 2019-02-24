const genSDK = (url, getToken, axios) => ({
  signUp: async (json) => {
    return await axios.post(`${url}/auth/signUp`, JSON.stringify(json));
  },
  signIn: async (json) => {
    return await axios.post(`${url}/auth/signIn`, JSON.stringify(json));
  },
  user: {
    me: async () => {
      return await axios.get(`${url}/users/me`, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
    get: async (userName) => {
      return await axios.get(`${url}/users/${userName}`, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
    follow: async (userName) => {
      return await axios.post(`${url}/users/${userName}/follow`, {}, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
    feed: {
      list: async (userName) => {
        return await axios.get(`${url}/users/${userName}/feed`, { headers: { Authorization: `Bearer ${getToken()}` } });
      },
    },
  },
  comment: {
    create: async (collectionId, message) => {
      return await axios.post(`${url}/comments`, {
        collectionId,
        message,
      }, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
    list: async (collectionId) => {
      return await axios.get(`${url}/collections/${collectionId}/comments`, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
  },
  collection: {
    create: async (form) => {
      return await axios.post(
        `${url}/collections`,
        form,
        { headers: { Authorization: `Bearer ${getToken()}` } }
      );
    },
    get: async (collectionId) => {
      return await axios.get(`${url}/collections/${collectionId}`, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
    list: async () => {
      return await axios.get(`${url}/collections`, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
    delete: async (collectionId) => {
      return await axios.delete(`${url}/collections/${collectionId}`, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
  },
  article: {
    create: async (collectionId, form) => {
      return await axios.post(
        `${url}/collections/${collectionId}/articles`,
        form,
        { headers: { Authorization: `Bearer ${getToken()}` } }
      );
    },
    generate_presigned_url: async (collectionId, key) => {
      return await axios.post(
        `${url}/collections/${collectionId}/articles-presigned`,
        key,
        { headers: { Authorization: `Bearer ${getToken()}` } }
      );
    },
    list: async (collectionId) => {
      return await axios.get(`${url}/collections/${collectionId}/articles`, { headers: { Authorization: `Bearer ${getToken()}` } });
    },
  },
});

module.exports = {
  genSDK,
};
