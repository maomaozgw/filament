import { ModelFilament } from '@/models/api';
import { search } from '@/services/warehouse/Filament';
import {
    ActionType,
    PageContainer, ProDescriptionsItemProps,
    ProTable
} from '@ant-design/pro-components';
import { useModel } from '@umijs/max';
import { Button, ColorPicker, Divider, message } from 'antd';
import React, { useRef, useState } from 'react';
import CreateForm from './components/CreateForm';
import ActionForm from './components/StockActionForm';


const Filament: React.FC<unknown> = () => {
    // 全量入库
    const [stockInModalVisible, handleStockInModalVisible] = useState<boolean>(false);
    // 快速库操作
    const [stockActionVisible, handleStockActionModalVisible] = useState<boolean>(false);
    const [stockActionName, setActionName] = useState<"stock-in" | "stock-out">("stock-in");
    const [currentFilament, setCurrentFilament] = useState<ModelFilament>({});
    const { initialState: { brands, colors, types } } = useModel("@@initialState");

    const actionRef = useRef<ActionType>();
    const inlineProcessProps = {
        editable: false // 不可编辑
    }
    const columns: ProDescriptionsItemProps<ModelFilament>[] = [
        {
            title: '类型',
            tip: '耗材的类型',
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
            min: 0,
        },
        {
            title: '价格',
            dataIndex: 'price',
            valueType: 'money', // TODO 并且乘以100传递给后端，从后端获取的数据也需要/100
            hideInTable: true,
            hideInSearch: true,
            precision: 2,
            min: 0,
        },
        {
            title: '更新时间',
            hideInForm: true,
            dataIndex: 'updated_at',
            hideInSearch: true,
            valueType: 'dateTime',
        },
        {
            title: '操作',
            dataIndex: 'option',
            valueType: 'option',
            render: (_, record) => (
                <>
                    <a
                        onClick={() => {
                            handleStockActionModalVisible(true);
                            setActionName("stock-in");
                            setCurrentFilament(record);
                        }}
                    >
                        入库
                    </a>
                    <Divider type="vertical" />
                    <a onClick={() => {
                        handleStockActionModalVisible(true);
                        setActionName("stock-out");
                        setCurrentFilament(record);
                    }}>出库</a>
                </>
            ),
        },
    ];
    return (
        <PageContainer
            header={{
                title: '耗材详情',
            }}
        >
            <ProTable<ModelFilament>
                headerTitle="查询耗材"
                rowKey='id'
                actionRef={actionRef}
                search={{
                    labelWidth: 100,
                }}
                columns={columns}
                toolBarRender={() => [
                    <Button key="1" type="primary" onClick={() => handleStockInModalVisible(true)}>
                        入库
                    </Button>,
                ]}
                request={async (params) => {
                    const { code, message: msg, data, pager: { total } } = await search(params);
                    if (code != 0) {
                        message.error(msg);
                        return { data: [], success: false };
                    }
                    return { data: data, success: true, total: total ?? 0 };
                }}
            />
            <CreateForm modalVisible={stockInModalVisible} onFinished={function (changed: boolean): void {
                handleStockInModalVisible(false);
                if (changed) {
                    actionRef.current?.reloadAndRest?.();
                }
            }}
            />
            <ActionForm
                modalVisible={stockActionVisible}
                filament={currentFilament}
                stockActionType={stockActionName}
                onFinish={(changed) => {
                    handleStockActionModalVisible(false);
                    if (changed) {
                        actionRef.current?.reload?.();
                    }
                }}
            />
        </PageContainer>
    )
}

export default Filament;