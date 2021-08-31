import TemplateComponent from './template-component'
import { render } from 'setupTests'

const setUp = () => {
  return render(<TemplateComponent/>)
}

describe('TemplateComponent component', () => {
  it('should render successfully', () => {
    const c = setUp()
    expect(c.container).toBeInTheDocument()
  })
})

