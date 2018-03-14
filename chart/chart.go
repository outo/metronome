package chart

import (
   "github.com/outo/metronome/track"
   "fmt"
   "path/filepath"
   "html/template"
   "os"
   "math"
   "strconv"
   "io/ioutil"
)

type TemplateData struct {
   EarliestStartNanos,
   LatestEndNanos int64
   Charts []*Chart
}

type Chart struct {
   Id              string
   Description     template.HTML
   Groups          []*Group
   Items           []*Item
   itemIdGenerator func() string
   CustomTimes     []int64
   CodeBase        string
}

type Group struct {
   Id      string
   Order   string
   Content template.HTML
   DefaultItemClassName,
   DefaultItemType string
}

func (g Group) WithContent(content string) Group {
   g.Content = template.HTML(content)
   return g
}

type Item struct {
   Id,
   Type,
   Group,
   Content,
   Title string
   StartNanos,
   EndNanos int64
   ClassName string
}

type TemplateExecutionMeta struct {
   track.ExecutionId
   Description string
}

type idGenerator func() string

func newIdGenerator() idGenerator {
   id := -1
   return func() string {
      id++
      return strconv.Itoa(id)
   }
}

type AllMetas map[track.ExecutionId]Meta

type Meta struct {
   Order       int
   Description string
   Groups      map[string]Group
   CodeBase    string
}

func extractTemplateData(allMetas AllMetas, eventsChannel <-chan track.Event) (data TemplateData) {

   data.EarliestStartNanos = int64(math.MaxInt64)
   data.LatestEndNanos = int64(math.MinInt64)

   // prepare charts placeholders
   data.Charts = make([]*Chart, len(allMetas))
   for executionId, meta := range allMetas {

      codeBase := meta.CodeBase
      if codeBase == "" {
         codeBase = string(executionId)
      }

      data.Charts[meta.Order] = &Chart{
         Id:              string(executionId),
         Description:     template.HTML(meta.Description),
         itemIdGenerator: newIdGenerator(),
         CodeBase:        codeBase,
      }
   }

   // process events channel
   for event := range eventsChannel {
      var chart *Chart
      for _, chart = range data.Charts {
         if chart.Id == string(event.ExecutionId) {
            break
         }
      }
      if chart == nil {
         break
      }

      if event.Category == "custom_time" {
         chart.CustomTimes = append(chart.CustomTimes, event.TimeSinceStart.Nanoseconds())
         continue
      }

      group := allMetas[track.ExecutionId(chart.Id)].Groups[event.Category]

      startNanos := event.TimeSinceStart.Nanoseconds()
      endNanos := startNanos + event.Duration.Nanoseconds()

      if startNanos < data.EarliestStartNanos {
         data.EarliestStartNanos = startNanos
      }

      if endNanos > data.LatestEndNanos {
         data.LatestEndNanos = endNanos
      }

      itemClassName := event.Category
      if event.ChartItemClassName != "" {
         itemClassName = event.ChartItemClassName
      } else if group.DefaultItemClassName != "" {
         itemClassName = group.DefaultItemClassName
      }

      itemType := group.DefaultItemType
      if event.ChartItemType != "" {
         itemType = event.ChartItemType
      }

      item := &Item{
         Id:         chart.itemIdGenerator(),
         Type:       itemType,
         Group:      group.Id,
         Content:    event.Content,
         Title:      event.Description,
         StartNanos: startNanos,
         EndNanos:   endNanos,
         ClassName:  itemClassName,
      }
      chart.Items = append(chart.Items, item)
   }

   data.EarliestStartNanos -= (data.LatestEndNanos - data.EarliestStartNanos) / 100
   data.LatestEndNanos += (data.LatestEndNanos - data.EarliestStartNanos) / 100

   // prepare groups
   for _, chart := range data.Charts {
      for _, group := range allMetas[track.ExecutionId(chart.Id)].Groups {
         if group.Content != "" {
            chart.Groups = append(chart.Groups, &Group{
               Id:      group.Id,
               Order:   group.Id,
               Content: group.Content,
            })
         }
      }
   }

   return
}

func generate(name, templatePath, outputPath string, data interface{}) (err error) {

   fmt.Println("generating:", name)

   templateAbsolutePath, err := filepath.Abs(templatePath)
   if err != nil {
      return
   }
   fmt.Println("Using chart template:", templateAbsolutePath)

   htmlTemplate, err := template.ParseFiles(templateAbsolutePath)
   if err != nil {
      return
   }

   chartAbsolutePath, err := filepath.Abs(outputPath)
   if err != nil {
      return
   }
   fmt.Println("Chart path:", chartAbsolutePath)

   file, err := os.Create(chartAbsolutePath)
   if err != nil {
      return
   }
   defer file.Sync()
   defer file.Close()

   err = htmlTemplate.Execute(file, data)
   if err != nil {
      return
   }

   return
}

func AllEvents(allMetas AllMetas, eventsChannel <-chan track.Event) (err error) {

   templateData := extractTemplateData(allMetas, eventsChannel)

   if err = generate(
      "javascript",
      "./templates/metronome-cacophony.js.template",
      "./docs/metronome-cacophony.js",
      templateData); err != nil {
      return
   }

   if err = generate(
      "accordion html",
      "./templates/chart-fragment.html.template",
      "./docs/chart-fragment.html",
      templateData); err != nil {
      return
   }

   bts, err := ioutil.ReadFile("./docs/chart-fragment.html")
   if err != nil {
      return
   }

   if err = generate(
      "page html",
      "./templates/page.html.template",
      "./docs/page.html",
      struct {
         Html template.HTML
      }{
         Html: template.HTML(bts),
      }); err != nil {
      return
   }

   return
}
