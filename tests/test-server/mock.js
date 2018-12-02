const mock = require('../mock-server');
const server = new mock.MockServer();

server.contract.get('/users/me', (req, res) => {
  res.send({
    id: '1',
    user_name: 'me',
  });
});

server.contract.get('/projects', (req, res) => {
  res.send([
    {
      id: '1',
      title: 'Project Meow',
      description: 'ぞうの卵はおいしいぞう。ぞうの卵はおいしいぞう。ぞうの卵はおいしいぞう。',
      media: [
        'document',
        'picture',
        'movie',
      ],
      cover: {
        sort: 'solid',
        color: 'teal darken-2',
      }
    },
    {
      id: '2',
      title: 'Piyo-piyo',
      description: '',
      media: [
        'document',
      ],
      cover: {
        sort: 'solid',
        color: 'orange lighten-2',
      }
    },
  ]);
});

server.contract.get('/projects/:projectId', (req, res) => {
  res.send({
    id: req.params.projectId,
    title: 'Project Meow',
    description: 'ぞうの卵はおいしいぞう。ぞうの卵はおいしいぞう。ぞうの卵はおいしいぞう。',
    owner: '1',
    media: [
      'document',
      'picture',
      'movie',
    ],
    cover: {
      sort: 'solid',
      color: 'teal darken-2',
    },
    collections: [
      {
        id: '1',
        name: 'My Collection',
        items: [
          {
            id: '1',
            type: 'share',
            entity: {
              source: 'ogp',
              ogp: {
                image: 'https://avatars2.githubusercontent.com/u/1870091?s=400&amp;v=4',
                title: 'myuon/pyxis',
                description: 'Pyxis. Contribute to myuon/pyxis development by creating an account on GitHub.',
                url: 'https://github.com/myuon/pyxis',
              }
            }
          },
          {
            id: '2',
            type: 'share',
            entity: {
              source: 'twitter',
              ogp: {
                image: 'https://pbs.twimg.com/media/DtTpA7TU4AAQNx9.jpg:large',
                title: 'みょん on Twitter',
                description: '“様子',
                url: 'https://twitter.com/myuon_myon/status/1068735246793756673'
              }
            }
          },
        ]
      }
    ],
    comments: [
      {
        id: '1',
        owner: {
          id: '1',
          user_name: 'me',
        },
        created_at: '2018/12/01 20:30:45',
        message: 'へいへいほー\nYes\nにゃんぽよ',
      },
      {
        id: '2',
        owner: {
          id: '2',
          user_name: 'he',
        },
        created_at: '2018/12/01 20:40:45',
        message: 'てすぽよ',
      },
      {
        id: '3',
        owner: {
          id: '3',
          user_name: 'she',
        },
        created_at: '2018/12/01 20:45:45',
        message: 'ぽよい',
      },
      {
        id: '4',
        owner: {
          id: '1',
          user_name: 'me',
        },
        created_at: '2018/12/01 21:50:55',
        message: 'へい',
      },
    ]
  });
});

module.exports = {
  server
};
