// const { Pact } = require('@pact-foundation/pact');

describe('The API', () => {
  it('should be trivially true', () => {
    expect('a').toEqual('a');
  });

  it('shoud not be trivially true', () => {
    expect('a').not.toEqual('b');
  });
});
