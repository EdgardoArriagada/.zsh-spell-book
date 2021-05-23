let g:singleQuoteFirst = ["'", '"']
vnoremap <silent> ' :<c-u>normal! \<esc>:<c-u>call PerfectlySelectString(g:singleQuoteFirst)<cr>
onoremap <silent> ' :<c-u>normal! :<c-u>call PerfectlySelectString(g:singleQuoteFirst)<cr>

let g:doubleQuoteFirst = ['"', "'"]
vnoremap <silent> " :<c-u>normal! \<esc>:<c-u>call PerfectlySelectString(g:doubleQuoteFirst)<cr>
onoremap <silent> " :<c-u>normal! :<c-u>call PerfectlySelectString(g:doubleQuoteFirst)<cr>

func! PerfectlySelectString(quotes)
  let g:savedPos = getpos('.')
  let g:savedColumn = g:savedPos[2]

  func! s:restoreValues()
    :call s:restorePos(g:savedPos)
    :call s:pressEscape()
  endfunc

  func! s:didMove()
    return col('.') != g:savedColumn
  endfunc


  " If did select between cursor
  for l:quote in a:quotes
    execute "normal! vi".l:quote

    if s:didSelectBetween(g:savedColumn)
      return
    endif

    :call s:restoreValues()
  endfor

  for l:quote in a:quotes
    " If did select goind forward
    execute "normal! vi".l:quote

    if s:didSelect() || s:didMove()
      return
    endif

    :call s:restoreValues()
  endfor

  for l:quote in a:quotes
    " If did select goind backwards
    execute "normal! F".l:quote."vi".l:quote

    if s:didSelect() || s:didMove()
      return
    endif

    :call s:restoreValues()
  endfor
endfunc

func! s:restorePos(position)
  :call setpos('.', a:position)
endfunc

func! s:pressEscape()
  execute "normal! \<esc>"
endfunc

func! s:didSelectBetween(middle)
  let l:left = col('v')
  let l:right = col('.')
  return l:left != l:right && l:left -1 <= a:middle && a:middle <= l:right +1
endfunc

func! s:didSelect()
  return col('.') != col('v')
endfunc

