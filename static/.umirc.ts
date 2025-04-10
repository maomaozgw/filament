import { defineConfig } from '@umijs/max';

export default defineConfig({
  antd: {},
  access: {},
  model: {},
  initialState: {},
  request: {},
  proxy: {
    '/api': {
      'target': "http://127.0.0.1:8080/",
      'changeOrigin': true,
    }
  },
  layout: {
    title: '3D打印耗材管理',
  },
  routes: [
    {
      path: '/',
      redirect: '/warehouse/dashboard',
    },
    {
      name: '库存',
      path: '/warehouse',
      routes: [
        {
          name: '耗材状态',
          path: '/warehouse/dashboard',
          component: './Warehouse/Statistic',
        },
        {
          name: '耗材详情',
          path: '/warehouse/list',
          component: './Warehouse/Filament',
        },
        {
          name: '耗材记录',
          path: '/warehouse/records',
          component: './Warehouse/Record',
        }
      ],
    },
   
    {
      name: '元数据管理',
      path: '/metadata',
      routes: [
        {
          name: '类型管理',
          path: '/metadata/type',
          component: './MetaData/Type',
        },
        {
          name: '品牌管理',
          path: '/metadata/brand',
          component: './MetaData/Brand',
        },
        {
          name: '颜色管理',
          path: '/metadata/color',
          component: './MetaData/Color',
        },
      ]
    },
  ],
  npmClient: 'pnpm',
});

