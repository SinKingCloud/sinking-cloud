import {Body} from "@/layouts/components";
import {NamePath} from "rc-field-form/es/interface";
import {Form} from "antd";


export default () => {
    const [form] = Form.useForm();
    const test = () => {
        form?.getFieldValue("test");
        form?.getFieldValue("test2" as NamePath)
    }
    return (
        <Body breadCrumb={false}>
            test
        </Body>
    );
};