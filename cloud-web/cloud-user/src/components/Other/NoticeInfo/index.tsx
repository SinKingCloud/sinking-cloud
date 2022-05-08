import React from 'react'
import {Modal} from 'antd'

class NoticeInfo extends React.Component {

  constructor(props: any) {
    super(props);//指定this，可以使用this
    this.state = {visible: props?.props || false}
  }

  componentDidMount() {
    console.log("初始化")
  }

  render() {
    return (
      <Modal title={"测试"} key={"noticeInfo"} width={350} destroyOnClose={true}
             forceRender={true}
        // @ts-ignore
             visible={this.state?.visible || false}>
      </Modal>
    )
  }
}

export default NoticeInfo
