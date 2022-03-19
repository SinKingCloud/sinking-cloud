import {DefaultFooter} from '@ant-design/pro-layout';
import Settings from '@/../config/defaultSettings';


export default () => {
  const defaultMessage = Settings.footer;

  const currentYear = new Date().getFullYear();

  return (
    <DefaultFooter
      copyright={`${currentYear} ${defaultMessage}`}
      links={[]}
    />
  );
};
