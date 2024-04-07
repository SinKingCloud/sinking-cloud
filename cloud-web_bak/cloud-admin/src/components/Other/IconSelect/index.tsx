import {Component} from 'react'
import {Button, Input, Popover} from "antd"
import styles from './style.less'
import * as allIcons from "@ant-design/icons";

import Icon from './icons'

export default class SelectIcon extends Component {

  static getDerivedStateFromProps(nextProps: any) {
    // Should be a controlled component.
    if ('value' in nextProps) {
      return {
        ...(nextProps.value || {}),
      };
    }
    return null;
  }

  constructor(props: any) {
    super(props);

    const value = props.value || '';
    this.state = {
      value,
    };
  }

  handleChangeAction = (value: any) => {
    if ('value' in this.props) {
      this.setState({value});
    }
    this.triggerChange(value);
  }

  triggerChange = (value: any) => {
    // @ts-ignore
    const onChange = this.props.onChange;
    if (onChange) {
      onChange(value);
    }
  }

  render() {
    // @ts-ignore
    const {value} = this.state;

    const icon_box = (
      <div className={styles.icon_inner}>
        <div className={styles.icon_list}>
          {Object.keys(allIcons).map((icon) => {
            if (allIcons[icon] != undefined && allIcons[icon] != null && allIcons[icon] != "" && typeof allIcons[icon] == "object") {
              return <Button key={icon}
                             icon={<Icon icon={icon}/>}
                             onClick={() => this.handleChangeAction(icon)}/>;
            }
            return null;
          })}
        </div>
      </div>
    );

    return (
      <div>
        <Popover
          content={icon_box}
          trigger="click"
        >
          <Input type="primary"
                 placeholder='请选择图标'
                 onChange={this.handleChangeAction}
                 value={value}
          />
        </Popover>
      </div>
    )
  }
}
