const express = require('express');

class MockServer {
  constructor () {
    this.app = express();
    this.server = null;

    this.contract = new express.Router();
  }

  start () {
    this.app.use(function(req, res, next) {
      res.header("Access-Control-Allow-Origin", "*");
      res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
      next();
    });

    this.app.use('/', this.contract);

    this.server = this.app.listen(5000, () => {
      console.log('listening on port 5000...');
    });
  }

  shutdown () {
    this.server.close();
  }
};

module.exports = {
  MockServer,
};
