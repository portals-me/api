const { Pact } = require('@pact-foundation/pact');
const path = require('path');
const collections = require('@/sdk/collection');

/*
const provider = new Pact({
  port: 8989,
  log: path.resolve(__dirname, 'logs', 'example.log'),
  dir: path.resolve(__dirname, 'pacts'),
  spec: 2,
  cors: true,
  pactfileWriteMode: 'update',
  consumer: 'Consumer',
  provider: 'Provider',
});

beforeAll(() => provider.setup());
afterEach(() => provider.verify());
afterAll(() => provider.finalize());
*/

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
