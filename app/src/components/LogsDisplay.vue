<template>
  <section>
        <b-field grouped group-multiline>
            <div v-for="(column, index) in columnsTemplate" 
                :key="index"
                class="control">
                <b-checkbox v-model="column.visible">
                    {{ column.name }}
                </b-checkbox>
            </div>
        </b-field>
  <b-table
    v-if="logs.length > 0"
    :data="logs"
    :default-sort="['time', 'desc']"
    :per-page="settings.logsPerPage"
    detailed
    paginated
    pagination-simple
  >
    <template slot-scope="props">
      <b-table-column v-for="(column, index) in columnsTemplate"
        :key="index"
        :label="column.name"
        :field="column.name"
        :visible="column.visible"
        sortable
      >{{ props.row[column.name] }}</b-table-column>
    </template>

    <template slot="detail" slot-scope="props">
      <span v-html="formatFullMessage(props.row)"/>
    </template>
  </b-table>
  <div v-else class="notification has-text-centered">
    <b-icon icon="package-variant" size="is-large"></b-icon>
    <p>
      No log messages received&hellip; yet!
    </p>
  </div>
  </section>
</template>

<script>
export default {
  name: "logs-display",
  props: ["logs", "settings"],
  data() {
    return {
      todayDateString: new Date().toLocaleDateString(),
      columnsTemplate: [
        { name: "time", visible: true },
        { name: "source", visible: true },
        { name: "level", visible: true },
        { name: "msg", visible: true }
      ]
    };
  },
  methods: {
    formatFullMessage(obj) {
      return "<pre>" + JSON.stringify(obj, null, "  ") + "</pre>";
    },
    severityClass(severity) {
      if (severity === 0) {
        // Emergicy
        return "is-dark";
      } else if (severity < 4) {
        // Alert, Critical, Error
        return "is-danger";
      } else if (severity === 4) {
        // Warning
        return "is-warning";
      } else if (severity === 7) {
        // Debug
        return "is-info";
      }

      return "is-light";
    }
  },
  watch: {
    logs: function() {
      var fields = [];
      for (const log of this.logs) {
        for (var prop in log) {
          if (fields.indexOf(prop) === -1) {
            fields.push(prop);
            var found = false;
            for (var col of this.columnsTemplate) {
              if (col.name == prop) {
                found = true;
                break;
              }
            }
            if (!found) {
              this.columnsTemplate.push({
                name: prop,
                visible: false
              });
            }
          }
        }
      }
      console.log(fields);
    }
  }
};
</script>
