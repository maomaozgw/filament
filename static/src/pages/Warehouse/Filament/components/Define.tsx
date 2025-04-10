import { ModelFilament } from "@/models/api";
import { InitialStateType } from "@/models/global";
import { ProDescriptionsItemProps } from "@ant-design/pro-components";
import { useModel } from "@umijs/max";
import { ColorPicker } from "antd";

export const filamentColumns = (extraProps: ProDescriptionsItemProps<ModelFilament>[]): ProDescriptionsItemProps<ModelFilament>[] => {
    const { initialState } = useModel("@@initialState");
    const { brands, colors, types } = initialState as InitialStateType;
    const result: ProDescriptionsItemProps<ModelFilament>[] = [
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
            title: '类型',
            tip: '耗材的类型',
            dataIndex: 'type_id',
            render: (_, record) => {
                if (record.type?.name == undefined) {
                    return <div>{record.type?.major}-{record.type?.minor}</div>;
                } else {
                    return <div>{record.type?.name}</div>;
                }

            },
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
            valueType: 'money',
            hideInSearch: true,
            precision: 2,
            min: 0,
            render: (_, record) => {
                return <div>{((record.price ?? 0) / 100).toFixed(2)}</div>;
            }
        },
        {
            title: '入库时间',
            hideInForm: true,
            dataIndex: 'created_at',
            hideInSearch: true,
            valueType: 'dateTime',
        },
    ];
    result.push(...extraProps);
    return result;
}