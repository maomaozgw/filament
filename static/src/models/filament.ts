
// export interface PageSearchRequest {
//     page: number
//     size: number
//     kind?: string
//     type?: string
//     type_id?: number
//     brand?: string
//     brand_id?: number
//     color?: string
//     color_id?: number
// }

export interface PageSearch {
    pageSize?: number;
    current?: number;
    keyword?: string;
}

export type PageSearchRequest<T> = T & PageSearch;