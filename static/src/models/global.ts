// 全局共享数据示例
import { DEFAULT_NAME } from '@/constants';
import { useState } from 'react';
import { ModelBrand, ModelColor, ModelType } from './api';

const useUser = () => {
  const [name, setName] = useState<string>(DEFAULT_NAME);
  return {
    name,
    setName,
  };
};

export default useUser;

export interface InitialStateType {
  brands: ModelBrand[]
  colors: ModelColor[]
  types: ModelType[]
}