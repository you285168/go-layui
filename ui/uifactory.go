package ui

type UIFactory interface {
	AddElem(e HtmlElem) UIFactory
	AddButton(text string, event OnButtonClick, assginto **UIButton) UIFactory
	AddLabel(text string) UIFactory
	AddRow() UIFactory
	AddText(text string, assginto **UIText) UIFactory
	AddRadio(selindex int, options []string, assginto **UIRadio) UIFactory
	AddSelect(selindex int, options []string, assginto **UISelect) UIFactory
	AddTextArea(prompt, text string, assginto **UITextArea) UIFactory
	AddSpace() UIFactory
	AddHref(text, url string) UIFactory
	AddLegend(text string) UIFactory
	AddEditer(prompt, text string, password bool, assginto **UIEditer) UIFactory
	AddTimePicker(format, displaytype string, val int64, assginto **UITimePicker) UIFactory
	AddUpload(text string, onupload UploadFile, assginto **UIUpload) UIFactory
	AddCheckBox(text string, checked bool, assginto **UICheckBox) UIFactory
	AddMergely(lf, rf string, onf OnGetFile, assginto **UIMergely) UIFactory
	AddStaticTable(header []string, data [][]string) *UITable
	AddTable(header []string, gd OnTableGetData) *UITable
}

type uiFactory struct {
	UIFactory
}

func NewFactory(f UIFactory) *uiFactory {
	return &uiFactory{UIFactory: f}
}

func (f *uiFactory) AddStaticTable(header []string, data [][]string) *UITable {
	return NewStaticTable(header, data)
}

func (f *uiFactory) AddTable(header []string, gd OnTableGetData) *UITable {
	return NewTable(header, gd)
}

func (f *uiFactory) AddButton(text string, event OnButtonClick, assginto **UIButton) UIFactory {
	b := NewButton(text, event)
	if assginto != nil {
		*assginto = b
	}
	return f.UIFactory.AddElem(b)
}

func (f *uiFactory) AddMergely(lf, rf string, onf OnGetFile, assginto **UIMergely) UIFactory {
	b := NewMergely(lf, rf, onf)
	if assginto != nil {
		*assginto = b
	}
	return f.UIFactory.AddElem(b)
}

func (f *uiFactory) AddCheckBox(text string, checked bool, assginto **UICheckBox) UIFactory {
	b := NewCheckBox(text, checked)
	if assginto != nil {
		*assginto = b
	}
	return f.UIFactory.AddElem(b)
}

func (f *uiFactory) AddUpload(text string, onupload UploadFile, assginto **UIUpload) UIFactory {
	b := NewUpload(text, onupload)
	if assginto != nil {
		*assginto = b
	}
	return f.UIFactory.AddElem(b)
}

func (f *uiFactory) AddLabel(text string) UIFactory {
	return f.UIFactory.AddElem(NewLabel(text))
}

func (f *uiFactory) AddRow() UIFactory {
	r := NewRow()
	f.UIFactory.AddElem(r)
	return r
}

func (f *uiFactory) AddText(text string, assginto **UIText) UIFactory {
	t := NewText(text)
	if assginto != nil {
		*assginto = t
	}
	return f.UIFactory.AddElem(t)
}

func (f *uiFactory) AddRadio(selindex int, options []string, assginto **UIRadio) UIFactory {
	r := NewRadio(selindex, options)
	if assginto != nil {
		*assginto = r
	}
	return f.UIFactory.AddElem(r)
}

func (f *uiFactory) AddSelect(selindex int, options []string, assginto **UISelect) UIFactory {
	s := NewSelect(selindex, options)
	if assginto != nil {
		*assginto = s
	}
	return f.UIFactory.AddElem(s)
}

func (f *uiFactory) AddTextArea(prompt, text string, assginto **UITextArea) UIFactory {
	ta := NewTextArea(prompt, text)
	if assginto != nil {
		*assginto = ta
	}
	return f.UIFactory.AddElem(ta)
}

func (f *uiFactory) AddSpace() UIFactory {
	return f.UIFactory.AddElem(NewSpace())
}

func (f *uiFactory) AddHref(text, url string) UIFactory {
	return f.UIFactory.AddElem(NewHref(text, url))
}

func (f *uiFactory) AddLegend(text string) UIFactory {
	return f.UIFactory.AddElem(NewLegend(text))
}

func (f *uiFactory) AddEditer(prompt, text string, password bool, assginto **UIEditer) UIFactory {
	e := NewEditer(prompt, text, password)
	if assginto != nil {
		*assginto = e
	}
	return f.UIFactory.AddElem(e)
}

func (f *uiFactory) AddTimePicker(format, displaytype string, val int64, assginto **UITimePicker) UIFactory {
	tp := NewTimePicker(format, displaytype, val)
	if assginto != nil {
		*assginto = tp
	}
	return f.UIFactory.AddElem(tp)
}
