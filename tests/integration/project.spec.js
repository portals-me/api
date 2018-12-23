const axios = require('axios');
const api = require('../../src/api');
const jwt = require('jsonwebtoken');
const fs = require('fs');

const testJWT = jwt.sign({
  id: 'user##test-user'
}, fs.readFileSync('./token/jwtES256.key', 'utf8'), { algorithm: 'ES256' });
const sdk = api.genSDK('http://localhost:3000', testJWT, axios);

describe('Project', () => {
  let projectCount = 0;

  it('can list projects', async () => {
    const result = await sdk.project.list();
    expect(result.status).toBe(200);

    projectCount = result.data.length;
  });

  let projectId;

  it('can create a project', async () => {
    const result = await sdk.project.create({
      title: 'project title',
      description: 'project description',
    });

    expect(result.status).toBe(201);

    projectId = result.headers.location.split('/projects/')[1];
  });

  it('can get a project', async () => {
    const result = await sdk.project.get(projectId);

    expect(result.status).toBe(200);
  });

  it('should have one more project after creation', async () => {
    const result = await sdk.project.list();
    expect(result.data.length).toBe(projectCount + 1);
  });
});
