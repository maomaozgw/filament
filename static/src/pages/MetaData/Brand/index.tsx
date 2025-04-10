import { ModelBrand } from "@/models/api";
import { ActionType, PageContainer, ProDescriptionsItemProps, ProTable, RequestData } from "@ant-design/pro-components";
import { request, useModel } from "@umijs/max";
import { Form } from "antd";
import { useRef, useState } from "react";
import EditForm from "./Components/EditForm";
import { InitialStateType } from "@/models/global";


async function queryBrands(params): Promise<Partial<RequestData<ModelBrand>>> {
    const { current, pageSize, ...filters } = params;
    const { data, pager } = await request("/api/v1/meta-data/brands", {
        method: "GET",
        params: {
            page: current,
            size: pageSize,
            ...filters,
        },
    });
    return {
        data: data,
        success: true,
        total: pager.total,
    }
}

const Brand: React.FC = () => {
    const actionRef = useRef<ActionType>();
    const [updateModalVisible, handleUpdateModalVisible] = useState<boolean>(false);
    const [currentItem, setCurrentItem] = useState<ModelBrand>({});
    const [updateForm] = Form.useForm<ModelBrand>();
    const columns: ProDescriptionsItemProps<ModelBrand>[] = [
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
    ]
    return (
        <PageContainer
            header={{
                title: "品牌管理",
            }}
        >
            <ProTable<ModelBrand>
                rowKey="id"
                actionRef={actionRef}
                columns={columns}
                request={queryBrands}
            />
            <EditForm
                currentItem={currentItem}
                actionRef={actionRef}
                updateModalVisible={updateModalVisible}
                handleUpdateModalVisible={handleUpdateModalVisible}
            ></EditForm>
        </PageContainer>
    )
}

export default Brand;