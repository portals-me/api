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
    expect(result.data.length).toBeGreaterThanOrEqual(0);

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

  describe('With a project', () => {
    it('should have one more project after creation', async () => {
      const result = await sdk.project.list();
      expect(result.data.length).toBe(projectCount + 1);
    });

    let commentCount = 0;

    it('can list comments', async () => {
      const result = await sdk.comment.list(projectId);
      expect(result.status).toBe(200);
      expect(result.data.length).toBeGreaterThanOrEqual(0);

      commentCount = result.data.length;
    });

    let commentIndex = null;

    it('can create a comment', async () => {
      const result = await sdk.comment.create(projectId, 'this is a message.');
      expect(result.status).toBe(201);

      commentIndex = parseInt(result.headers.location.split('/')[4], 10);
      expect(commentIndex).toBeGreaterThanOrEqual(0);
    });

    it('should be consistent with comment_count', async () => {
      const project = (await sdk.project.get(projectId)).data;

      expect(project.comment_count).toBe(commentIndex + 1);
    });
  });
});
