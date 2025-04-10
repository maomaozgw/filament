import { ApiItemResponseModelFilament, ApiListResponseModelFilament, ModelFilament } from "@/models/api";
import { PageSearchRequest } from "@/models/filament";
import { request } from "@umijs/max";
import { message } from "antd";

export async function search(params: PageSearchRequest<ModelFilament>): Promise<ApiListResponseModelFilament> {
    const { current, pageSize, ...filters } = params;
    const res = await request("/api/v1/warehouse/filaments", {
        method: "GET",
        params: {
            page: current,
            size: pageSize,
            ...filters,
        },
    });
    return res;

}

export async function stockIn(value: ModelFilament): Promise<boolean> {
    const hide = message.loading('入库中');
    try {
        value.price = Math.floor((value.price ?? 0)*100);
        const resp = await request("/api/v1/warehouse/filaments", {
            method: "POST",
            data: value,
        });
        hide();
        message.success('入库成功');
        return true;
    } catch (error) {
        // TODO 优化错误信息展现
        hide();
        message.error(`入库失败${error?.message ?? ""}，请重试`);
        return false;
    }
}

export async function stockOut(value: ModelFilament): Promise<boolean>{
    const hide = message.loading('出库中');
    try {
        await request(`/api/v1/warehouse/filaments/${value.id}`, {
            method: "delete",
            data: {
                quantity: value.quantity,
            },
        });
        hide();
        message.success('出库成功');
        return true;
    } catch (error) {
        hide();
        message.error('出库失败，请重试');
        return false;
    }
}