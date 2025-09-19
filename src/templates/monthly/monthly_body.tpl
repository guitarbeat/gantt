% Setup category palette for this month
\SetupDefaultCategoryPalette

{{- template "calendar_table.tpl" dict "Cfg" .Cfg "Body" .Body -}}


% Single full-width Notes area (replaces previous two-column layout)
\parbox{\textwidth}{
  \myUnderline{Notes}
  \vbox to \dimexpr\textheight-\pagetotal-\myLenLineHeightButLine\relax {%
    \leaders\vbox to \myLenLineHeightButLine{\vfil\hrule width \linewidth height \myLenLineThicknessDefault}\vskip \dimexpr\textheight-\pagetotal-\myLenLineHeightButLine\relax
  }%
}
