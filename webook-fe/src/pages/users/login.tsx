import React, { useState } from 'react';
import { Button, Form, Input, Space, Typography, Divider } from 'antd';
import axios from "@/axios/axios";
import Link from "next/link";
import router from "next/router";

const { Text } = Typography;

const LoginForm: React.FC = () => {
    const [loading, setLoading] = useState(false);

    const onFinish = async (values: any) => {
        try {
            setLoading(true);
            const res = await axios.post("/users/login", values);
            if (res.status !== 200) {
                alert(res.statusText);
                return;
            }
            alert(res.data);
            router.push('/users/profile');
        } catch (err) {
            alert("登录失败：" + err);
        } finally {
            setLoading(false);
        }
    };

    const onFinishFailed = () => {
        alert("输入有误");
    };

    return (
        <Form
            name="login_form"
            labelCol={{ span: 6 }}
            wrapperCol={{ span: 14 }}
            style={{ maxWidth: 500, margin: '0 auto', paddingTop: 50 }}
            onFinish={onFinish}
            onFinishFailed={onFinishFailed}
            autoComplete="off"
        >
            <Form.Item
                label="邮箱"
                name="email"
                rules={[{ required: true, message: '请输入邮箱' }]}
            >
                <Input placeholder="请输入邮箱" />
            </Form.Item>

            <Form.Item
                label="密码"
                name="password"
                rules={[{ required: true, message: '请输入密码' }]}
            >
                <Input.Password placeholder="请输入密码" />
            </Form.Item>

            <Form.Item wrapperCol={{ offset: 6, span: 14 }}>
                <Button type="primary" htmlType="submit" loading={loading} block>
                    登录
                </Button>
            </Form.Item>

            <Form.Item wrapperCol={{ offset: 6, span: 14 }}>
                <Divider plain>其他方式登录</Divider>
                <Space direction="vertical" style={{ width: '100%' }}>
                    <Link href="/users/login_sms">
                        <Button block>使用手机号登录</Button>
                    </Link>
                    <Link href="/users/login_wechat">
                        <Button block>微信扫码登录</Button>
                    </Link>
                </Space>
            </Form.Item>

            <Form.Item wrapperCol={{ offset: 6, span: 14 }}>
                <Text>还没有账号？</Text>
                &nbsp;
                <Link href="/users/signup">
                    <Text type="secondary" underline>去注册</Text>
                </Link>
            </Form.Item>
        </Form>
    );
};

export default LoginForm;
