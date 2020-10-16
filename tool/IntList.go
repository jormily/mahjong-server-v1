package tool


type IntList []int

func NewIntList(len,cap int) IntList  {
	this := make(IntList,len,cap)
	return this
}

func (this *IntList)Insert(value int,args...int) bool {
	index := len(*this)+1
	if len(args) == 1 && args[0] > 0 {
		index = args[0]
	}
	if index > len(*this) {
		index = len(*this)+1
	}
	if index <= 0 {
		index = 1
	}

	index = index - 1
	temp := append([]int{}, (*this)[index:]...)
	*this = append(append((*this)[:index], value), temp...)
	return true
}

func (this *IntList)RemoveByValue(value int,args...int){
	cnt := 1
	if len(args) == 1 && args[0] > 0 {
		cnt = args[0]
	}

	for i:=len(*this)-1;i>=0;i--{
		if (*this)[i] == value {
			(*this) = append((*this)[:i], (*this)[i+1:]...)
			cnt--
			if cnt == 0 {
				return
			}
		}
	}
}

func (this *IntList)IndexOf(value int) int {
	for i:=0;i<len(*this);i++{
		if (*this)[i] == value {
			return i+1
		}
	}
	return 0
}

func (this *IntList)Remove(index int) bool {
	if !(index > 0 && index <= len(*this)) {
		return false
	}

	index = index - 1
	*this = append((*this)[:index], (*this)[index+1:]...)
	return true
}

func (this *IntList)Clone() IntList {
	return append([]int{},*this...)
}

func (this *IntList)Map() map[int]int {
	m := make(map[int]int)
	for _,v := range *this {
		if _,ok := m[v];ok {
			m[v] = m[v] + 1
		}
	}
	return m
}

func (this *IntList)Len() int {
	return len(*this)
}
