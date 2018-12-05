export const api = (endpoint, axios) => {
  return {
    getMe: async () => {
      const result = await axios.get(`${endpoint}/users/me`);
      return result.data;
    },
  };
};
