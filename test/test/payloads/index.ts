import {
  generateRandomEmail,
  generateRandomName,
  generateRandomPassword,
  generateTestObjects,
} from "../helper/generator";
import { ActivityTypes_E } from "../schemas/activity.schema";
import { HeightUnitTypes, PreferenceTypes, WeightUnitTypes } from "../schemas/user.schema";
import type { Payload } from "./payload.types";

const positiveUser = {
  preference: "WEIGHT",
  weightUnit: "LBS",
  heightUnit: "INCH",
  weight: 75,
  height: 167,
  name: generateRandomName(),
  imageUri: ""
}

const positiveFile = {
  path: "img-8kb.jpeg"
}

const hugeFile = {
  path: "img-155kb.png"
}

const positiveRegisterEmail = {
  email: generateRandomEmail(),
  password: generateRandomPassword(8, 32),
};

const positiveRegisterAnotherAccount = {
  email: "another" + generateRandomEmail(),
  password: generateRandomPassword(8, 32),
};

const positiveActivity = {
  activityType: "JumpRope",
  doneAt: new Date(Date.now() + 15 * 60_000).toISOString(),
  durationInMinutes: 15
};

const positiveActivityUpdate = {
  activityType: "Running",
  doneAt: new Date(Date.now() + 45 * 60_000).toISOString(),
  durationInMinutes: 45
};

export const payloads: Record<string, Payload> = {
  userPreference: {
    positive: positiveUser,
    negative: generateTestObjects({
      preference: {
        type: "string",
        enum: PreferenceTypes,
        notNull: true
      },
      weightUnit: {
        type: "string",
        enum: WeightUnitTypes,
        notNull: true
      },
      heightUnit: {
        type: "string",
        enum: HeightUnitTypes,
        notNull: true
      },
      weight: {
        type: "number",
        min: 10,
        max: 1000,
        notNull: true
      },
      height: {
        type: "number",
        min: 3,
        max: 250,
        notNull: true
      },
      name: {
        type: "string",
        minLength: 2,
        maxLength: 60,
        notNull: false
      },
      imageUri: {
        type: "string",
        notNull: false
      }
    }, positiveUser),
  },
  fileHuge: {
    positive: hugeFile,
    negative: generateTestObjects({
      path: {
        notNull: true,
        type: "string"
      }
    }, hugeFile),
  },
  filePas: {
    positive: positiveFile,
    negative: generateTestObjects({
      path: {
        notNull: true,
        type: "string"
      }
    }, positiveFile),
  },
  register: {
    positive: positiveRegisterEmail,
    negative: generateTestObjects({
      email: {
        notNull: true,
        isEmail: true,
        type: "string"
      },
      password: {
        notNull: true,
        minLength: 8,
        maxLength: 32,
        type: "string"
      }
    }, positiveRegisterEmail),
  },
  registerAnother: {
    positive: positiveRegisterAnotherAccount,
    negative: []
  },
  login: {
    positive: positiveRegisterEmail,
    negative: generateTestObjects({
      email: {
        notNull: true,
        isEmail: true,
        type: "string",
      },
      password: {
        notNull: true,
        minLength: 8,
        maxLength: 32,
        type: "string"
      }
    }, positiveRegisterEmail),
  },
  loginInvalidPass: {
    positive: {
      phone: positiveRegisterEmail.email,
      password: "asdf",
    },
    negative: [],
  },
  // Activity
  // POST
  activity: {
    positive: positiveActivity,
    negative: generateTestObjects({
      activityType: {
        type: "string",
        enum: ActivityTypes_E,
        notNull: true
      },
      doneAt: {
        type: "string",
        notNull: true
      },
      durationInMinutes: {
        type: "number",
        notNull: true,
        min: 1,
        max: 199
      }
    }, positiveActivity),
  },
  updateActivity: {
    positive: positiveActivityUpdate,
    negative: []
  }
};
