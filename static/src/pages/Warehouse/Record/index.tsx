import { ModelFilament, ModelRecord } from '@/models/api';
import { PageSearchRequest } from '@/models/filament';
import { ActionType, ModalForm, PageContainer, ProDescriptionsItemProps, ProForm, ProFormDigit, ProFormMoney, ProFormText, ProTable } from '@ant-design/pro-components';
import { request, useModel, useRequest } from '@umijs/max';
import { Button, ColorPicker, Divider, Form, message, Tag } from 'antd';
import React, { useRef, useState } from 'react';

const queryRecords = async (params) => {
    const { current, pageSize, ...filter } = params;
    const req_params: PageSearchRequest = {
        page: current,
        size: pageSize,
        ...filter,
    }

    const res = await request("/api/v1/warehouse/records", {
        method: "GET",
        params: req_params,
    });
    return res;
}

const Record: React.FC<unknown> = () => {
    const actionRef = useRef<ActionType>();
    const { initialState: { brands, colors, types } } = useModel("@@initialState");
    const inlineProcessProps = {
        editable: false // 不可编辑
    }
    const columns: ProDescriptionsItemProps<ModelRecord>[] = [
        {
            title: "操作",
            dataIndex: "kind",
            valueEnum: {
                "stock-in": {
                    text: "入库"
                },
                "stock-out": {
                    text: "出库"
                },
                "stock-take": {
                    text: "盘点"
                },
            },
        },
        {
            title: '类型',
            dataIndex: 'type_id',
            render: (_, record) => (
                <div>{record.type?.major}-{record.type?.minor}</div>
            ),
            editable: false,
            request: async (value) => {
                return types.map((type) => {
                    return {
                        value: type.id,
                        label: `${type.major}-${type.minor}`,
                    }
                });
            }
        },
        {
            title: '品牌',
            dataIndex: 'brand_id',
            render: (_, record) => (
                <div>{record.brand?.name}</div>
            ),
            valueType: 'text',
            request: async (value) => {
                return brands.map((brand) => {
                    return {
                        value: brand.id,
                        label: brand.name,
                    }
                });
            }
        },
        {
            title: '颜色',
            dataIndex: 'color_id',
            valueType: 'text',
            render: (_, record) => (
                <ColorPicker defaultValue={record.color?.rgba} disabled showText={() => record.color?.name} />
            ),
            request: async (value) => {
                return colors.map((color) => {
                    return {
                        value: color.id,
                        label: color.name,
                    }
                });
            }
        },
        {
            title: '数量',
            dataIndex: 'quantity',
            valueType: 'digit',
            hideInSearch: true,
        },
        {
            title: '单价',
            dataIndex: 'price',
            valueType: 'money',
            hideInSearch: true,
            render: (_, record) => {
                if (record.kind == "stock-in") {
                    return <div >¥{(record.price ?? 0) / 100}</div>
                }
                return <></>
            }
        },
        {
            title: '时间',
            hideInForm: true,
            dataIndex: 'created_at',
            valueType: 'dateTime',
        },
    ];
    return (
        <PageContainer
            header={{
                title: '耗材记录',
            }}
        >
            <ProTable<ModelFilament>
                headerTitle="查询记录"
                rowKey='id'
                actionRef={actionRef}
                search={{
                    labelWidth: 100,
                }}
                columns={columns}
                request={async (params, sorter, filter) => {
                    const { code, message, data, pager: { total } } = await queryRecords(params);
                    if (code != 0) {
                        message.error(message);
                        return { data: [], success: false };
                    }
                    return { data: data, success: true, total: total ?? 0 };
                }}
            />
        </PageContainer>
    )
}

export default Record;