import { ModelType } from "@/models/api";
import { InitialStateType } from "@/models/global";
import { ActionType, PageContainer, ProDescriptionsItemProps, ProTable, RequestData, RequestOptionsType } from "@ant-design/pro-components";
import { request, useModel } from "@umijs/max";
import { Form } from "antd";
import { useRef, useState } from "react";


async function qeuryTypes(param): Promise<Partial<RequestData<ModelType>>> {
    const { current, pageSize, ...filters } = param;
    const { data, pager } = await request("/api/v1/meta-data/types", {
        method: "GET",
        params: {
            page: current,
            size: pageSize,
            ...filters,
        }
    });
    return {
        success: true,
        data: data,
        total: pager.total,
    }
}



const Type: React.FC = () => {
    const actionRef = useRef<ActionType>();
    const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
    const [currentItem, setCurrentItem] = useState<ModelType>({});
    const [updateForm] = Form.useForm<ModelType>();
    const { initialState } = useModel("@@initialState");
    const { types } = initialState as InitialStateType;
    const columns: ProDescriptionsItemProps<ModelType>[] = [
        {
            title: "id",
            dataIndex: "id",
            valueType: "digit",
            hideInTable: true,
            hideInForm: true,
            hideInSearch: true,
        },
        {
            title: "名称",
            dataIndex: "name",
            valueType: "text",
            hideInSearch: true,
        },
        {
            title: "主要分类",
            dataIndex: "major",
            valueType: "text",
            request: async (value) => {
                const result: RequestOptionsType[] = [];
                const dataSet = new Set<string>();
                types.forEach((element) => {
                    if (dataSet.has(element.major ?? "")) {
                        return;
                    }
                    dataSet.add(element.major ?? "");
                    result.push({
                        value: element.major,
                        label: element.major,
                    });
                });
                return result;
            }
        },
        {
            title: "次要分类",
            dataIndex: "minor",
            valueType: "text",
            request: async (value) => {
                const result: RequestOptionsType[] = [];
                const dataSet = new Set<string>();
                types.forEach((element) => {
                    if (dataSet.has(element.minor ?? "")) {
                        return;
                    }
                    dataSet.add(element.minor ?? "");
                    result.push({
                        value: element.minor,
                        label: element.minor,
                    });
                });
                return result;
            }
        },
        {
            title: "更新时间",
            dataIndex: "updated_at",
            valueType: "dateTime",
            hideInForm: true,
            hideInSearch: true,
        },
        {
            title: "操作",
            dataIndex: "option",
            valueType: "option",
            render: (_, record) => (
                <>
                    <a onClick={() => {
                        handleUpdateModalVisible(true);
                        setCurrentItem(record);
                    }}  >编辑</a>
                </>
            )
        }
    ];
    return (
        <PageContainer
            header={{
                title: "耗材类型管理",
            }}
        >
            <ProTable<ModelType>
                rowKey="id"
                actionRef={actionRef}
                columns={columns}
                request={qeuryTypes}
            />
        </PageContainer>
    )
}

export default Type;