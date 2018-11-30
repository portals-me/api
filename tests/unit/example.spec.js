const collections = require('@/sdk/collection');

describe('The API', () => {
  describe('Collections', () => {
    it('should returns only own collections', async () => {
      const results = await collections.api.list('me');
      results.forEach(item => {
        expect(item.owner).toEqual('me');
      });
    });
  });
});
