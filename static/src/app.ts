// 运行时配置

import { RunTimeLayoutConfig } from "@umijs/max";
import type { RequestConfig } from 'umi';
import { MetaData } from '@/services/metadata';
import logo from './assets/logo.png'
import { ModelBrand, ModelColor } from "./models/api";
import { InitialStateType } from "./models/global";

// 全局初始化数据配置，用于 Layout 用户信息和权限初始化
// 更多信息见文档：https://umijs.org/docs/api/runtime-config#getinitialstate
export async function getInitialState(): Promise<InitialStateType> {
  const { brands, colors, types } = await MetaData.getMetaData() ?? {};
  return { brands: brands ?? [], colors: colors ?? [], types: types ?? [] };
}

export const layout: RunTimeLayoutConfig = () => {
  return {
    title: '3D打印耗材管理',
    logo: logo,

    fixedHeader: true,
    menu: {
      locale: true,
    },
  };
};

export const request: RequestConfig = {
  timeout: 1000,
  // other axios options you want
  errorConfig: {
    errorHandler() {
    },
    errorThrower() {
    }
  },
  requestInterceptors: [],
  responseInterceptors: []
};