export const api = {
  list: async (userId) => {
    return [
      { id: '1', owner: userId },
      { id: '2', owner: userId },
    ];
  },
};
