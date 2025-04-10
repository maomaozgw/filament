import { ApiItemResponseV1MetaData } from "@/models/api";
import { request } from "@umijs/max";

export async function getMetaData() {
    try {
        const { data } = await request<ApiItemResponseV1MetaData>("/api/v1/meta-data");
        const { brands, colors, types } = data ?? {};
        return {
            brands,
            colors,
            types,
        };
    } catch (error) {
        console.error(error);
    }
}