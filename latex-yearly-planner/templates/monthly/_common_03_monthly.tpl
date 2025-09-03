{{- template "monthTabularV2.tpl" dict "Cfg" .Cfg "Body" .Body -}}
\smallskip

\parbox{\myLenTwoCol}{
  \myUnderline{Notes}
  \vbox to \dimexpr\textheight-\pagetotal-\myLenLineHeightButLine\relax {%
    \leaders\vbox to \myLenLineHeightButLine{\vfil\hrule width \linewidth height \myLenLineThicknessDefault}\vskip \dimexpr\textheight-\pagetotal-\myLenLineHeightButLine\relax
  }%
}%
\hspace{\myLenTwoColSep}%
\parbox{\myLenTwoCol}{
  \myUnderline{Notes}
  \vbox to \dimexpr\textheight-\pagetotal-\myLenLineHeightButLine\relax {%
    \leaders\vbox to \myLenLineHeightButLine{\vfil\hrule width \linewidth height \myLenLineThicknessDefault}\vskip \dimexpr\textheight-\pagetotal-\myLenLineHeightButLine\relax
  }%
}
