import { request } from "@umijs/max";
import { message } from "antd";

export async function createImport<T>(kind: string, data: T): Promise<boolean> {
    const hide = message.loading('数据导入中');
    try {
        const payload = {
            kind:kind,
            data:  data,
        }
        const resp = await request("/api/v1/imports", {
            method: "POST",
            data: payload,
        });
        hide();
        message.success('数据导入成功');
        return true;
    } catch (error) {
        hide();
        message.error(`数据导入失败${error?.message ?? ""}，请重试`);
        return false;
    }
}