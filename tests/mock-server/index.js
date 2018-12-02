const express = require('express');

class MockServer {
  constructor () {
    this.app = express();
    this.server = null;

    this.contract = new express.Router();
  }

  start () {
    this.app.use('/', this.contract);
    this.server = this.app.listen(5000, () => {
      console.log('listening on port 5000...');
    });
  }

  shutdown () {
    this.server.close();
  }
};

const s = new MockServer();
s.contract.get('/hoge', (req, res, next) => {
  res.send({
    result: 'OK',
  });
});
s.start();

setTimeout(() => {
  s.shutdown();
}, 10000);
