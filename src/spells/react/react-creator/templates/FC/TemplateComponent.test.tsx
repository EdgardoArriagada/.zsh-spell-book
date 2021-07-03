import TemplateComponent from './template-component'
import { render } from 'setupTests'

describe('TemplateComponent component', () => {
  it('should render successfully', () => {
    const c = render(<TemplateComponent/>)
    expect(c.container).toBeInTheDocument()
  })
})

