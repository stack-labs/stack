import {Card, Col, Input, Row, Space} from 'antd';
import React, {FC, useState, useMemo, useRef } from 'react';

import {GridContent, PageHeaderWrapper} from '@ant-design/pro-layout';
import {connect, Dispatch} from 'umi';
import {RouteChildrenProps} from 'react-router';
import {CallState} from './model';
import Forms from '../../../components/Forms';


const {TextArea} = Input;


interface CallProps extends RouteChildrenProps {
  dispatch: Dispatch;
  callService: CallState;
  loading: boolean;
}

interface SelectOption {
  value:string,
  label:string
  request?:string
}



const Call: FC<CallProps> = ({ dispatch ,callService }) => {
  const [addOpt,setAddOpt] = useState<SelectOption[]>([])
  const [endpointOpt,setEndpointOpt] = useState<SelectOption[]>([])
  const { services, callResult } = callService
  let formRef = useRef<Forms | null>(null)
  const serviceOptions = useMemo(():SelectOption[]=>{
    let list: SelectOption[] = []
    if(services && services.length){
      list = services.map((item:any)=>({value:item.name,label:item.name}))
    }
    return list
  },[services])

  console.log(callService)

  const getName = (values:any):object=>{
    let obj = {}
    if(values){
      values.forEach((item:any)=>{
        obj[item.name] = ''
        if(item.values){
          obj[item.name] = getName(item.values)
        }

      })
    }
    return obj
  }


  const formOption = ():object[]=>[
    {
      formType:'select',
      name:'Service',
      options:serviceOptions,
      placeholder:'Service',
      rules:[{ required: true, message: 'Please input Service!' }],
      onChange:(val:any)=>{
        if(val){
          let selectItem: any = services.find((v:any)=>v.name===val)
          let addOpt: SelectOption[] = selectItem.nodes.map((v:any)=>({value:v.address,label:v.address}))
          let endpointOpt: SelectOption[] = selectItem.endpoints.map((v:any)=>({value:v.name,label:v.name,request:v.request}))
          setAddOpt(addOpt)
          setEndpointOpt(endpointOpt)
          if(formRef.current){
            formRef.current.setFieldsValue({Address:'',Endpoint:'',request:''})
          }
        }

      }
    },
    {
      formType:'select',
      name:'Address',
      options:addOpt,
      placeholder:'Address'
    },
    {
      formType:'select',
      name:'Endpoint',
      options:endpointOpt,
      placeholder:'Endpoint',
      rules:[{ required: true, message: 'Please input Endpoint!' }],
      onChange:(val:any)=>{
        if(val){
          let node:any = endpointOpt.find((v:any)=>v.value===val)
          let request:object = getName(node.request && node.request.values)
          formRef.current && formRef.current.setFieldsValue({request:JSON.stringify(request ,null,2)})
        }

      }
    },
    {
      formType:'textarea',
      name:'request',
      rows:12,
      rules:[{ required: true, message: 'Please input request!' }],
    },
    {
      formType:'button',
      label:'call',
      htmlType:"submit",
      onClick:()=>{
        if(formRef.current){
          formRef.current.validateFields((validate:any)=>{
            console.log(validate)
              if(validate.errorFields) return
              dispatch({type:'callService/callServicer',payload:{...validate,request:JSON.parse(validate.request)}})



          })
        }
      }
    }
  ]



  return (
    <PageHeaderWrapper>
      <GridContent>
        <Row gutter={24}>
          <Col lg={10} md={24}>
            <Card bordered={false} style={{marginBottom: 24}}>
              <Space style={{width: "100%"}} size="large" direction="vertical">
                <Forms formOpt={{className:'form'}} ref={formRef}  options={formOption()} />
              </Space>
            </Card>
          </Col>
          <Col lg={14} md={24}>
            <Card>
              <TextArea rows={23}  readOnly></TextArea>
            </Card>
          </Col>
        </Row>
      </GridContent>
    </PageHeaderWrapper>
  );
}

export default connect(
  ({
     loading,
     callService,
   }: {
    loading: { effects: { [key: string]: boolean } };
    callService: CallState;
  }) => ({
    callService,
    loading: loading.effects['call/fetch'],
  }),
)(Call);
